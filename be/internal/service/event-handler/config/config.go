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
	UserService            common.ServiceConfig  `yaml:"user_service"`
	WalletService          common.ServiceConfig  `yaml:"wallet_service"`
	General                General               `yaml:"general"`
	NewUserConsumer        common.ConsumerConfig `yaml:"new_user_consumer"`
	NewTransactionConsumer common.ConsumerConfig `yaml:"new_transaction_consumer"`
}

type General struct {
	InitialBalance int64 `yaml:"initial_balance"`
}
