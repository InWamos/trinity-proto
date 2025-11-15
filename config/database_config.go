package config

import "github.com/spf13/viper"

type DatabaseConfig struct {
	Address      string `mapstructure:"DATABASE_ADDRESS"`
	Port         int    `mapstructure:"DATABASE_PORT"`
	DatabaseName string `mapstructure:"DATABASE_NAME"`
	DatabaseUser string `mapstructure:"DATABASE_USER"`
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	var databaseConfig DatabaseConfig
	if err := viper.Unmarshal(&databaseConfig); err != nil {
		return nil, err
	}
	return &databaseConfig, nil
}
