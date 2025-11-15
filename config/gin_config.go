package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	BindAddress  string `mapstructure:"GIN_ADDRESS"`
	Port         int    `mapstructure:"GIN_PORT"`
	Mode         string `mapstructure:"GIN_MODE"`
	TrustedProxy string `mapstrucrure:"GIN_TRUSTED_PROXY"`
}

func NewServerConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}
	return &serverConfig, nil
}
