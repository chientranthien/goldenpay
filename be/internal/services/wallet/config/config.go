package config

import (
	"github.com/chientranthien/goldenpay/internal/common"
)

var (
	defaultConfig = &Config{}
)

func init() {
	common.GetDefaultConfig(defaultConfig)
}

func Get() Config {
	if defaultConfig == nil {
		common.L().Info("defaultConfigIsNil")
		return Config{}
	}
	return *defaultConfig
}

type Config struct {
	DB                     common.DBConfig       `yaml:"db"`
	WalletService          common.ServiceConfig  `yaml:"wallet_service"`
	NewTransactionProducer common.ProducerConfig `yaml:"new_transaction_producer"`
}
