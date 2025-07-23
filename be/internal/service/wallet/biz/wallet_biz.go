package biz

import (
	"fmt"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm/clause"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/service/wallet/dao"
)

type WalletBiz struct {
	dao                          *dao.WalletDao
	producer                     sarama.SyncProducer
	newTransactionProducerConfig common.ProducerConfig
}

func NewWalletBiz(
	dao *dao.WalletDao,
	producer sarama.SyncProducer,
	newTransactionProducerConfig common.ProducerConfig,
) *WalletBiz {
	return &WalletBiz{
		dao:                          dao,
		producer:                     producer,
		newTransactionProducerConfig: newTransactionProducerConfig,
	}
}

func (b WalletBiz) Transfer(req *proto.TransferReq) (*proto.TransferResp, error) {
	var err error
	dao, err := b.dao.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction, err=%v", err)
	}
	defer dao.CommitOrRollback(err == nil)

	w, err := dao.GetByUserIDForUpdate(req.FromUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallet, err=%v", err)
	}

	if w.Balance < req.Amount {
		return nil, status.Errorf(codes.InvalidArgument, "insufficient balance")
	}

	if w.Status != common.WalletStatusActive {
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet status")
	}

	w.Balance -= req.Amount
	err = dao.Update(w)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update wallet")
	}

	trans := &proto.Transaction{
		FromUser:   req.FromUser,
		ToUser:     req.ToUser,
		FromWallet: w.Id,
		Amount:     req.Amount,
		Status:     common.TransactionStatusPending,
		Version:    common.FirstVersion,
		Ctime:      common.NowMillis(),
		Mtime:      common.NowMillis(),
	}

	err = dao.InsertTransaction(trans)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert transaction, err=%v", err)
	}

	e := &proto.NewTransactionEvent{
		TransactionId: trans.Id,
		EventTime:     trans.Ctime,
		FromUser:      req.FromUser,
		ToUser:        req.ToUser,
	}

	msg := &sarama.ProducerMessage{
		Topic: b.newTransactionProducerConfig.Topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", req.FromUser)),
		Value: sarama.ByteEncoder(common.ToJsonIgnoreErr(e)),
	}

	// ignore produceErr, so if there is any error, it will be handled by transaction_fixer cronjob
	partition, offset, produceErr := b.producer.SendMessage(msg)
	common.L().Infow(
		"sentNewTransactionEvent",
		"partition", partition,
		"offset", offset,
		"err", produceErr,
	)

	return &proto.TransferResp{TransactionId: trans.Id}, nil
}

func (b WalletBiz) Topup(req *proto.TopupReq) (*proto.TopupResp, error) {
	var err error
	dao, err := b.dao.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction, err=%v", err)
	}
	defer dao.CommitOrRollback(err == nil)

	w, err := dao.GetByUserIDForUpdate(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallet, err=%v", err)
	}

	if w.Status != common.WalletStatusActive {
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet status")
	}

	w.Balance += req.Amount
	err = dao.Update(w)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update wallet")
	}

	t := &proto.Topup{
		UserId:   req.UserId,
		WalletId: w.Id,
		Amount:   req.Amount,
		Status:   common.TopupStatusSuccess,
		Version:  common.FirstVersion,
		Ctime:    common.NowMillis(),
		Mtime:    common.NowMillis(),
	}

	err = dao.InsertTopup(t)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert topup, err=%v", err)
	}

	return &proto.TopupResp{TopupId: t.Id}, nil
}

func (b WalletBiz) Get(req *proto.GetWalletReq) (*proto.GetWalletResp, error) {
	w, err := b.dao.GetByUserID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get from DB, err=%v", err)
	}

	resp := &proto.GetWalletResp{Balance: w.Balance}

	return resp, nil
}

func (b WalletBiz) GetByUserID(userID uint64) (*proto.Wallet, error) {
	return b.dao.GetByUserID(userID)
}

func (b WalletBiz) Create(req *proto.CreateWalletReq) (*proto.CreateWalletResp, error) {
	w := &proto.Wallet{
		UserId:   req.UserId,
		Balance:  req.InitialBalance,
		Currency: common.GoldenDollar,
		Status:   common.WalletStatusActive,
		Version:  common.FirstVersion,
		Ctime:    common.NowMillis(),
		Mtime:    common.NowMillis(),
	}

	if err := b.dao.Insert(w); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert to DB, err=%v", err)
	}

	return &proto.CreateWalletResp{}, nil
}

func (b WalletBiz) ProcessTransfer(req *proto.ProcessTransferReq) (*proto.ProcessTransferResp, error) {
	var err error
	dao, err := b.dao.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction, err=%v", err)
	}
	defer dao.CommitOrRollback(err == nil)

	transaction, err := b.dao.GetTransactionForUpdate(req.TransactionId)
	if err != nil {
		return nil, err
	}

	if transaction.Status != common.TransactionStatusPending {
		return nil, status.Errorf(codes.Code(common.CodeProceeded.Id), "transaction is proceeded")
	}

	toWallet, err := dao.GetByUserIDForUpdate(transaction.ToUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallet, err=%v", err)
	}

	if b.shouldApproveTransaction(transaction, toWallet) {
		toWallet.Balance += transaction.Amount
		err = dao.Update(toWallet)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update wallet, err=%v", err)
		}
		transaction.ToWallet = toWallet.Id
		transaction.Status = common.TransactionStatusSuccess
		err = dao.UpdateTransaction(transaction)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update transaction, err=%v", err)
		}
	} else {
		fromWallet, err := dao.GetByUserIDForUpdate(transaction.FromUser)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get wallet, err=%v", err)
		}

		fromWallet.Balance += transaction.Amount
		err = dao.Update(fromWallet)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update wallet, err=%v", err)
		}

		transaction.ToWallet = toWallet.Id
		transaction.Status = common.TransactionStatusRejected
		err = dao.UpdateTransaction(transaction)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update transaction, err=%v", err)
		}
	}

	return &proto.ProcessTransferResp{}, nil
}

func (b WalletBiz) shouldApproveTransaction(t *proto.Transaction, w *proto.Wallet) bool {
	if w.Status != common.WalletStatusActive {
		return false
	}

	if t.FromWallet == t.ToWallet {
		return false
	}

	return true
}

func (b WalletBiz) GetUserTransactions(req *proto.GetUserTransactionsReq) (*proto.GetUserTransactionsResp, error) {
	req.Pagination = common.EnsurePagination(req.Pagination)
	fromUser := clause.Eq{
		Column: dao.TransactionColFromUser,
		Value:  req.GetCond().User.Eq,
	}
	toUser := clause.Eq{
		Column: dao.TransactionColToUser,
		Value:  req.GetCond().User.Eq,
	}
	var whereCond clause.Expression
	whereCond = clause.OrConditions{Exprs: []clause.Expression{fromUser, toUser}}

	if req.Pagination.Val != 0 {
		id := clause.Lte{
			Column: dao.TransactionColId,
			Value:  req.Pagination.Val,
		}

		whereCond = clause.AndConditions{Exprs: []clause.Expression{whereCond, id}}
	}

	limit := clause.Limit{
		Limit: common.Int(int(req.Pagination.Limit)),
	}

	order := clause.OrderBy{
		Columns: []clause.OrderByColumn{{Column: clause.Column{Name: dao.TransactionColId}, Desc: true}},
	}

	transs, err := b.dao.GetUserTransactions(whereCond, order, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user transactions")
	}

	nextPagination := common.NexPagination(req.Pagination, func(e interface{}) int64 {
		if trans, ok := e.(*proto.Transaction); ok {
			return int64(trans.Id)
		}

		return 0
	}, transs)

	if nextPagination.HasMore {
		transs = transs[:len(transs)-1]
	}
	return &proto.GetUserTransactionsResp{Transactions: transs, NextPagination: nextPagination}, nil
}
