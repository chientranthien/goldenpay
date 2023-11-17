package biz

import "github.com/chientranthien/goldenpay/internal/services/user/dao"

type WalletBiz struct {
	dao *dao.WalletDao
}

func NewWalletBiz(dao *dao.WalletDao) *WalletBiz {
	return &WalletBiz{dao: dao}
}

func (b WalletBiz) GetBy()  {

}

