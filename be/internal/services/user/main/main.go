package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/user/biz"
	"github.com/chientranthien/goldenpay/internal/services/user/config"
	"github.com/chientranthien/goldenpay/internal/services/user/controller"
	"github.com/chientranthien/goldenpay/internal/services/user/dao"
	"github.com/chientranthien/goldenpay/internal/services/user/service"
)

func main() {
	lis, err := net.Listen("tcp", config.Get().UserService.Addr)
	if err != nil {
		log.Fatalf("failed to listen, addr=%v, err=%v", config.Get().UserService.Addr, err)
	}
	server := grpc.NewServer()
	db, err := gorm.Open(
		mysql.Open(config.Get().DB.GetDSN()),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		log.Fatalf("failed to open db, conf=%v, err=%v", config.Get().DB, err)
	}
	dao := dao.NewUserDao(db)
	biz := biz.NewUserBiz(config.Get().JWT, dao)
	proto.RegisterUserServiceServer(
		server,
		service.NewService(
			controller.NewSignupController(biz),
			controller.NewLoginController(biz),
			controller.NewGetController(biz),
			controller.NewAuthzController(biz),
		),
	)

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve, err=%v", err)
	}
}
