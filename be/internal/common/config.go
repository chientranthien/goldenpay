package common

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

const (
	defaultConfigFile = "config.yaml"
)

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

type ProducerConfig struct {
	Addrs   []string `yaml:"addrs"`
	Version string   `yaml:"version"`
	Topic   string   `yaml:"topic"`
}

type ConsumerConfig struct {
	Addrs         []string `yaml:"addrs"`
	Version       string   `yaml:"version"`
	Topic         string `yaml:"topic"`
	ConsumerGroup string   `yaml:"consumer_group"`
}

type ServiceConfig struct {
	Addr string `yaml:"addr"`
}

func GetCurrentDir() string {
	_, file, _, _ := runtime.Caller(3)
	return filepath.Dir(file)
}

func GetDefaultConfigFile() string {
	return GetCurrentDir() + "/" + defaultConfigFile
}

func GetDefaultConfig(c any) {
	f, err := os.Open(GetDefaultConfigFile())
	if err != nil {
		log.Fatalf("failed to open config file, err=%v", err)
	}

	decoder := yaml.NewDecoder(f)
	decoder.Decode(c)
	if err != nil {
		log.Fatalf("failed to unmarshal config, err=%v")
	}
}
