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

	// Create DAOs
	channelDao := dao.NewChannelDao(db)
	messageDao := dao.NewMessageDao(db)

	// Create business logic services
	channelBiz := biz.NewChannelBiz(channelDao)
	messageBiz := biz.NewMessageBiz(
		messageDao,
		channelDao,
		common.NewKafkaProducer(config.Get().NewMessageProducer),
		config.Get().NewMessageProducer,
	)

	server := server.NewServer(
		config.Get().ChatService,
		controller.NewGetChannelController(channelBiz),
		controller.NewSendMessageController(messageBiz),
		controller.NewGetMessagesController(messageBiz),
		controller.NewCreateDirectMessageChannelController(channelBiz),
		controller.NewGetDirectMessagesController(channelBiz),
	)

	server.Serve()
}
