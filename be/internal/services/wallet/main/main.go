package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/IBM/sarama"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/proto"
	"github.com/chientranthien/goldenpay/internal/services/wallet/biz"
	"github.com/chientranthien/goldenpay/internal/services/wallet/config"
	"github.com/chientranthien/goldenpay/internal/services/wallet/controller"
	"github.com/chientranthien/goldenpay/internal/services/wallet/dao"
	"github.com/chientranthien/goldenpay/internal/services/wallet/service"
)

func main() {
	lis, err := net.Listen("tcp", config.Get().WalletService.Addr)
	if err != nil {
		log.Fatalf("failed to listen, addr=%v, err=%v", config.Get().WalletService.Addr, err)
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

	dao := dao.NewWalletDao(db)
	biz := biz.NewWalletBiz(dao, NewKafkaProducer(config.Get().Kafka), config.Get().Kafka)
	proto.RegisterWalletServiceServer(
		server,
		service.NewService(
			controller.NewTransferController(biz),
			controller.NewGetController(biz),
		),
	)

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve, err=%v", err)
	}
}

func NewKafkaProducer(conf common.KafkaConfig) sarama.SyncProducer {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		log.Fatalf("failed to parse kafka version, err=%v", err)
	}
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.Addrs, config)
	if err != nil {
		log.Fatalf("failed to new producer, err=%v", err)
	}

	return producer
}
