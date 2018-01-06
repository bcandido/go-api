package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/op/go-logging"
)

const MODULE = "cofing"
var log = logging.MustGetLogger(MODULE)

// Config stores the application-wide configurations
var Config appConfig

type appConfig struct {
	DB map[interface{}]interface{} `mapstructure:"db"`
	ServerPort int `mapstructure:"port"`
}

func LoadConfig(configPaths ...string) error {
	log.Info("load application configurations")

	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.SetDefault("port", 8080)
	config.AutomaticEnv()

	for _, path := range configPaths {
		config.AddConfigPath(path)
	}

	err := config.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	err = config.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}
