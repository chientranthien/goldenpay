package biz

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/dao"
)

type WalletBiz struct {
	walletDao   *dao.WalletDao
	producer    sarama.SyncProducer
	kafkaConfig common.KafkaConfig
}

func NewWalletBiz(
	dao *dao.WalletDao,
	producer sarama.SyncProducer,
	kafkaConfig common.KafkaConfig,
) *WalletBiz {
	return &WalletBiz{
		walletDao:   dao,
		producer:    producer,
		kafkaConfig: kafkaConfig,
	}
}

type NewTransactionEvent struct {
	TransactionId uint64
	EventTime     uint64
}

func (b WalletBiz) Transfer(req *proto.TransferReq) (*proto.TransferResp, error) {
	var err error
	dao := b.walletDao.Begin()
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
		FromUser: req.FromUser,
		ToUser:   req.ToUser,
		Amount:   req.Amount,
		Status:   common.TransactionStatusPending,
		Version:  common.FirstVersion,
		Ctime:    common.NowMillis(),
		Mtime:    common.NowMillis(),
	}

	err = dao.InsertTransaction(trans)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert transaction, err=%v", err)
	}

	e := &NewTransactionEvent{
		TransactionId: trans.Id,
		EventTime:     trans.Ctime,
	}

	msg := &sarama.ProducerMessage{
		Topic: b.kafkaConfig.Topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", req.FromUser)),
		Value: sarama.ByteEncoder(common.ToJsonIgnoreErr(e)),
	}

	// ignore produceErr, so if there is any error, it will be handled by transaction_fixer cronjob
	partition, offset, produceErr := b.producer.SendMessage(msg)
	log.Println("sent kafka message, partition=%v, offset=%v, err=%v", partition, offset, produceErr)

	return &proto.TransferResp{}, nil
}
