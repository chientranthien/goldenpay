package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/service/chat/biz"
	"github.com/chientranthien/goldenpay/internal/service/chat/config"
	"github.com/chientranthien/goldenpay/internal/service/chat/controller"
	"github.com/chientranthien/goldenpay/internal/service/chat/dao"
	"github.com/chientranthien/goldenpay/internal/service/chat/server"
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

	dao := dao.NewChatDao(db)
	biz := biz.NewChatBiz(
		dao,
		common.NewKafkaProducer(config.Get().NewMessageProducer),
		config.Get().NewMessageProducer,
		common.NewKafkaProducer(config.Get().PresenceUpdateProducer),
		config.Get().PresenceUpdateProducer,
		common.NewKafkaProducer(config.Get().MembershipEventProducer),
		config.Get().MembershipEventProducer,
	)
	server := server.NewServer(
		config.Get().ChatService,
		controller.NewChannelController(biz),
		controller.NewMessageController(biz),
		controller.NewPresenceController(biz),
	)

	server.Serve()
}
