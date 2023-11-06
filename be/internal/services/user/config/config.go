package config

import (
	"log"

	"github.com/go-sql-driver/mysql"

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
	DB          DBConfig             `yaml:"db"`
	UserService common.ServiceConfig `yaml:"user_service"`
	JWT         JWTConfig            `yaml:"jwt"`
}

type DBConfig struct {
	Addr   string `yaml:"addr"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	DBName string `yaml:"db_name"`
}

func (c DBConfig) GetDSN() string {
	conf := mysql.Config{
		User:                 c.User,
		Passwd:               c.Pass,
		Net:                  "tcp",
		Addr:                 c.Addr,
		DBName:               c.DBName,
		AllowNativePasswords: true,
	}
	return conf.FormatDSN()
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	DurationInMin int64  `yaml:"duration_in_min"`
}
