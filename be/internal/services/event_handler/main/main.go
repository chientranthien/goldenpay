package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"github.com/chientranthien/goldenpay/internal/common"
	"github.com/chientranthien/goldenpay/internal/services/event_handler/config"
	"github.com/chientranthien/goldenpay/internal/services/event_handler/controller"
	walletclient "github.com/chientranthien/goldenpay/internal/services/wallet/client"
)

func main() {
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(2)
	go func() {
		defer wg.Done()
		wClient := walletclient.NewWalletServiceClient(config.Get().WalletService.Addr)
		controller := controller.NewNewUserController(wClient)
		newUserConsumerGroup := newConsumerGroup(config.Get().NewUserConsumer)
		newUserConsumerGroup.Consume(ctx, []string{config.Get().NewUserConsumer.Topic}, controller)
	}()
	go func() {
		defer wg.Done()
		wClient := walletclient.NewWalletServiceClient(config.Get().WalletService.Addr)
		controller := controller.NewNewTransactionController(wClient)
		newTransactionConsumerGroup := newConsumerGroup(config.Get().NewTransactionConsumer)
		newTransactionConsumerGroup.Consume(ctx, []string{config.Get().NewTransactionConsumer.Topic}, controller)
	}()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-exitCh
	cancel()
	wg.Wait()

}

func newConsumerGroup(conf common.ConsumerConfig) sarama.ConsumerGroup {
	saramaConfig := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		common.L().Fatalw("parseKafkaVersionErr", "conf", saramaConfig, "err", err)
	}
	saramaConfig.Version = version
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Return.Successes = true

	group, err := sarama.NewConsumerGroup(conf.Addrs, conf.ConsumerGroup, saramaConfig)
	if err != nil {
		common.L().Fatalw("newConsumerGroupErr", "conf", conf, "err", err)
	}

	return group
}
