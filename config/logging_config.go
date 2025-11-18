package config

import "github.com/spf13/viper"

type LoggingConfig struct {
	Level string `mapstructure:"LOGGING_LEVEL"`
}

func NewLoggingConfig() (*LoggingConfig, error) {
	var loggingConfig LoggingConfig
	if err := viper.Unmarshal(&loggingConfig); err != nil {
		return nil, err
	}
	return &loggingConfig, nil
}
