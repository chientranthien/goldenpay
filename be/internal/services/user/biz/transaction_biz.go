package biz

import (
	"github.com/chientranthien/goldenpay/internal/services/user/dao"
)

type TransactionBiz struct {
	dao *dao.TransactionDao
}

func NewTransactionBiz(dao *dao.TransactionDao) *TransactionBiz {
	return &TransactionBiz{dao: dao}
}


