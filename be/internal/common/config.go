package common

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigFile = "config.yaml"
)

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

func GetDefaultConfig(c any){
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