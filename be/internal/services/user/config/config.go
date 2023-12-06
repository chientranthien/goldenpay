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
	DB              common.DBConfig       `yaml:"db"`
	UserService     common.ServiceConfig  `yaml:"user_service"`
	JWT             JWTConfig             `yaml:"jwt"`
	NewUserProducer common.ProducerConfig `yaml:"new_user_producer"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	DurationInMin int64  `yaml:"duration_in_min"`
}
