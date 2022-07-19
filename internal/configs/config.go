package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	LoggerConfig `mapstructure:"logger"`
}

type LoggerConfig struct {
	Level       int    `mapstructure:"level"`
	InfoLogFile string `mapstructure:"info_log_file"`
}

func LoadConfig() (config *Config, l *logrus.Logger, err error) {
	// Initialize properties config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	l = InitLogger(&config.LoggerConfig)
	return config, l, err
}
