package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
	"github.com/chientranthien/goldenpay/internal/services/wallet/config"
	"github.com/chientranthien/goldenpay/internal/services/wallet/controller"
	"github.com/chientranthien/goldenpay/internal/services/wallet/dao"
	"github.com/chientranthien/goldenpay/internal/services/wallet/server"
)

func main() {
	db, err := gorm.Open(
		mysql.Open(config.Get().DB.GetDSN()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		common.L().Fatalw("openDBErr", "conf", config.Get().DB, "err", err)
	}

	dao := dao.NewWalletDao(db)
	biz := biz.NewWalletBiz(
		dao,
		common.NewKafkaProducer(config.Get().NewTransactionProducer),
		config.Get().NewTransactionProducer,
	)
	server := server.NewServer(
		config.Get().WalletService,
		controller.NewTransferController(biz),
		controller.NewTopupController(biz),
		controller.NewGetController(biz),
		controller.NewCreateController(biz),
		controller.NewGetUserTransactionsController(biz),
	)
	server.Serve()
}
