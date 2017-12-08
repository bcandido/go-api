package app

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cofing")

// Config stores the application-wide configurations
var Config appConfig

type appConfig struct {
	DB map[interface{}]interface{} `mapstructure:"db"`
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
// Environment variables with the prefix "RESTFUL_" in their names are also read automatically.
func LoadConfig(configPaths ...string) error {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("yaml")
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
