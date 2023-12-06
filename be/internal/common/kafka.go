package common

import "github.com/IBM/sarama"

func NewKafkaProducer(conf ProducerConfig) sarama.SyncProducer {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(conf.Version)

	if err != nil {
		L().Fatalw("parseKafkaVersionErr", "conf", config, "err", err)
	}
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.Addrs, config)
	if err != nil {
		L().Fatalw("newKafkaProducerFata", "conf", conf, "err", err)
	}

	return producer
}
