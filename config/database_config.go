package config

import "github.com/spf13/viper"

type DatabaseConfig struct {
	Address          string `mapstructure:"DATABASE_ADDRESS"`
	Port             int    `mapstructure:"DATABASE_PORT"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode  string `mapstructure:"DATABASE_SSL_MODE"`
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("DATABASE_ADDRESS")
	_ = viper.BindEnv("DATABASE_PORT")
	_ = viper.BindEnv("DATABASE_NAME")
	_ = viper.BindEnv("DATABASE_USER")
	_ = viper.BindEnv("DATABASE_PASSWORD")
	_ = viper.BindEnv("DATABASE_SSL_MODE")

	var databaseConfig DatabaseConfig
	if err := viper.Unmarshal(&databaseConfig); err != nil {
		return nil, err
	}
	return &databaseConfig, nil
}
