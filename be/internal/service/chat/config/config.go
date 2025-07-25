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
	DB                 common.DBConfig       `yaml:"db"`
	ChatService        common.ServiceConfig  `yaml:"chat_service"`
	NewMessageProducer common.ProducerConfig `yaml:"new_message_producer"`
}
