package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/service/user/biz"
	"github.com/chientranthien/goldenpay/internal/service/user/config"
	"github.com/chientranthien/goldenpay/internal/service/user/controller"
	"github.com/chientranthien/goldenpay/internal/service/user/dao"
	"github.com/chientranthien/goldenpay/internal/service/user/server"
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

	dao := dao.NewUserDao(db)
	biz := biz.NewUserBiz(
		config.Get().JWT,
		dao,
		common.NewKafkaProducer(config.Get().NewUserProducer),
		config.Get().NewUserProducer,
	)
	server := server.NewServer(
		config.Get().UserService,
		controller.NewSignupController(biz),
		controller.NewLoginController(biz),
		controller.NewGetController(biz),
		controller.NewAuthzController(biz),
		controller.NewGetByEmailController(biz),
	)

	server.Serve()
}
