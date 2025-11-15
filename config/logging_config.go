package config

import "github.com/spf13/viper"

type LoggingConfig struct {
	Level string `mapstructure:"LOGGING_LEVEL"`
	Out   string `mapstructure:"LOGGING_OUT"`
}

func NewLoggingConfig() (*LoggingConfig, error) {
	var loggingConfig LoggingConfig
	if err := viper.Unmarshal(&loggingConfig); err != nil {
		return nil, err
	}
	return &loggingConfig, nil
}
