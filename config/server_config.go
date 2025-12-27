package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Environment   string `mapstructure:"SERVER_ENVIRONMENT"`
	BindAddress   string `mapstructure:"SERVER_ADDRESS"`
	Port          int    `mapstructure:"SERVER_PORT"`
	TrustedProxy  string `mapstructure:"SERVER_TRUSTED_PROXY"`
	AllowedOrigin string `mapstructure:"SERVER_ALLOWED_ORIGIN"`
}

func NewServerConfig() (*ServerConfig, error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("SERVER_ENVIRONMENT")
	_ = viper.BindEnv("SERVER_ADDRESS")
	_ = viper.BindEnv("SERVER_PORT")
	_ = viper.BindEnv("SERVER_TRUSTED_PROXY")
	_ = viper.BindEnv("SERVER_ALLOWED_ORIGIN")

	var serverConfig ServerConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}
	return &serverConfig, nil
}
