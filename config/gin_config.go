package config

import (
	"github.com/spf13/viper"
)

type GinConfig struct {
	BindAddress   string `mapstructure:"GIN_ADDRESS"`
	Port          int    `mapstructure:"GIN_PORT"`
	Mode          string `mapstructure:"GIN_MODE"`
	TrustedProxy  string `mapstructure:"GIN_TRUSTED_PROXY"`
	AllowedOrigin string `mapstructure:"GIN_ALLOWED_ORIGIN"`
}

func NewServerConfig() (*GinConfig, error) {
	var serverConfig GinConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}
	return &serverConfig, nil
}
