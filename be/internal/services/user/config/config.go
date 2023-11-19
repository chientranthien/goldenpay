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
	DB          common.DBConfig      `yaml:"db"`
	UserService common.ServiceConfig `yaml:"user_service"`
	JWT         JWTConfig            `yaml:"jwt"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	DurationInMin int64  `yaml:"duration_in_min"`
}
