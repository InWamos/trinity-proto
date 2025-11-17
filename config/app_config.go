package config

import "github.com/spf13/viper"

type AppConfig struct {
	LoggingConf    LoggingConfig  `mapstructure:",squash"`
	DatabaseConfig DatabaseConfig `mapstructure:",squash"`
	GinConfig      GinConfig      `mapstructure:",squash"`
}

func NewAppConfig() (*AppConfig, error) {
	viper.AutomaticEnv()
	var appConfig AppConfig

	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
