package config

import (
	"log"

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
		log.Println("default config is nil")
		return Config{}
	}
	return *defaultConfig
}

type Config struct {
	DB            common.DBConfig      `yaml:"db"`
	WalletService common.ServiceConfig `yaml:"wallet_service"`
	Kafka         common.KafkaConfig   `yaml:"kafka"`
}
