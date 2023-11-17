package dao

import (
	"gorm.io/gorm"

	"github.com/chientranthien/goldenpay/internal/proto"
)

type TransactionDao struct {
	db *gorm.DB
}

func NewTransactionDao(db *gorm.DB) *TransactionDao {
	return &TransactionDao{db: db}
}

func (d TransactionDao) getDB() *gorm.DB {
	return d.db.Table("transaction_tab")
}

func (d TransactionDao) Insert(trans *proto.Transaction) error {
	return d.getDB().Create(trans).Error
}

func (d *TransactionDao) Update(trans *proto.Transaction) error {
	return d.getDB().Updates(trans).Commit().Error
}

func (d *TransactionDao) Get(id uint64) (*proto.Transaction, error) {
	u := &proto.Transaction{
		Id: id,
	}

	if err := d.getDB().First(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}
