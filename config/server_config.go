package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	BindAddress   string `mapstructure:"SERVER_ADDRESS"`
	Port          int    `mapstructure:"SERVER_PORT"`
	TrustedProxy  string `mapstructure:"SERVER_TRUSTED_PROXY"`
	AllowedOrigin string `mapstructure:"SERVER_ALLOWED_ORIGIN"`
}

func NewServerConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}
	return &serverConfig, nil
}
