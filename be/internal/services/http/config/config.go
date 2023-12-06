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
	HttpService   common.ServiceConfig `yaml:"http_service"`
	UserService   common.ServiceConfig `yaml:"user_service"`
	WalletService common.ServiceConfig `yaml:"wallet_service"`
}
