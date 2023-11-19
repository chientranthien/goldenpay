package dao

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

func (d WalletDao) Insert(wallet *proto.Wallet) error {
	return d.getDB().Create(wallet).Error
}

func (d *WalletDao) Update(wallet *proto.Wallet) error {
	return d.getDB().Updates(wallet).Error
}


func (d WalletDao) Begin() *WalletDao {
	trx := d.getDB().Begin()
	return &WalletDao{db: trx}
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
		log.Printf("failed to get wallet from DB for update, err=%v", err)
		return nil, err
	}

	return w, nil
}

func (d WalletDao) InsertTransaction(trans *proto.Transaction) error {
	return d.getTransactionDB().Create(trans).Error
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
