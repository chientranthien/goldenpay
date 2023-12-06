package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
)

type WalletDao struct {
	db *gorm.DB
}

func NewWalletDao(db *gorm.DB) *WalletDao {
	return &WalletDao{db: db}
}

func (d WalletDao) getDB() *gorm.DB {
	return d.db.Table("wallet_tab")
}

func (d WalletDao) getTransactionDB() *gorm.DB {
	return d.db.Table("transaction_tab")
}

func (d WalletDao) getTopupDB() *gorm.DB {
	return d.db.Table("topup_tab")
}

func (d WalletDao) Insert(wallet *proto.Wallet) error {
	if err := d.getDB().Create(wallet).Error; err != nil {
		common.L().Errorw("insertErr", "wallet", wallet, "err", err)
		return err
	}

	return nil
}

func (d *WalletDao) Update(wallet *proto.Wallet) error {
	wallet.Mtime = common.NowMillis()
	wallet.Version++
	if err := d.getDB().Updates(wallet).Error; err != nil {
		common.L().Errorw("updateErr", "wallet", wallet, "err", err)
		return err
	}

	return nil
}

func (d WalletDao) Begin() (*WalletDao, error) {
	trx := d.getDB().Begin()
	if err := trx.Error; err != nil {
		common.L().Errorw("beginErr", "err", err)
		return nil, err
	}

	return &WalletDao{db: trx}, nil
}

func (d WalletDao) CommitOrRollback(commit bool) {
	if commit {
		d.getDB().Commit()
	} else {
		d.getDB().Rollback()
	}
}

func (d *WalletDao) GetByUserIDForUpdate(userID uint64) (*proto.Wallet, error) {
	w := &proto.Wallet{
		UserId: userID,
	}

	err := d.getDB().Clauses(clause.Locking{Strength: "UPDATE"}).First(w).Error
	if err != nil {
		common.L().Infow("getByUserIDForUpdateErr", "userID", userID, "err", err)
		return nil, err
	}

	return w, nil
}

func (d *WalletDao) Get(id uint64) (*proto.Wallet, error) {
	w := &proto.Wallet{
		Id: id,
	}

	if err := d.getDB().First(w).Error; err != nil {
		return nil, err
	}

	return w, nil
}

func (d *WalletDao) GetByUserID(userID uint64) (*proto.Wallet, error) {
	w := &proto.Wallet{
		UserId: userID,
	}

	err := d.getDB().First(w).Error
	if err != nil {
		common.L().Infow("getByUserIDeErr", "userID", userID, "err", err)
		return nil, err
	}

	return w, nil
}

func (d WalletDao) InsertTransaction(trans *proto.Transaction) error {
	return d.getTransactionDB().Create(trans).Error
}

func (d WalletDao) InsertTopup(t *proto.Topup) error {
	return d.getTopupDB().Create(t).Error
}

func (d *WalletDao) GetTransactionForUpdate(id uint64) (*proto.Transaction, error) {
	t := &proto.Transaction{
		Id: id,
	}

	err := d.getTransactionDB().Clauses(clause.Locking{Strength: "UPDATE"}).First(t).Error
	if err != nil {
		common.L().Infow("getTransactionForUpdate", "transactionID", id, "err", err)
		return nil, err
	}

	return t, nil
}

func (d WalletDao) UpdateTransaction(t *proto.Transaction) error {
	t.Mtime = common.NowMillis()
	t.Version++
	if err := d.getTransactionDB().Updates(t).Error; err != nil {
		common.L().Errorw("updateErr", "transaction", t, "err", err)
		return err
	}

	return nil
}

func (d WalletDao) GetUserTransactions(conds ...clause.Expression) ([]*proto.Transaction, error) {

	var transactions []*proto.Transaction
	if err := d.getTransactionDB().Clauses(conds...).Find(&transactions).Error; err != nil {
		common.L().Errorw("getUserTransactionsErr", "conds", conds, "err", err)
		return nil, err
	}

	return transactions, nil
}
