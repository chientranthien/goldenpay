package dao

import (
	"gorm.io/gorm"

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

func (d WalletDao) Insert(wallet *proto.Wallet) error {
	return d.getDB().Create(wallet).Error
}

func (d *WalletDao) Update(wallet *proto.Wallet) error {
	return d.getDB().Updates(wallet).Commit().Error
}

func (d *WalletDao) Get(id uint64) (*proto.Wallet, error) {
	u := &proto.Wallet{
		Id: id,
	}

	if err := d.getDB().First(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}
