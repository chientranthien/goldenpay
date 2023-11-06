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

func GetDefaultConfig() Config {
	if defaultConfig == nil {
		log.Println("default config is nil")
		return Config{}
	}
	return *defaultConfig
}

type Config struct {
	HttpService   common.ServiceConfig `yaml:"http_service"`
	UserService   common.ServiceConfig `yaml:"user_service"`
	WalletService common.ServiceConfig `yaml:"wallet_service"`
}
