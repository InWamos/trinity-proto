package config

import "github.com/spf13/viper"

type RedisConfig struct {
	Host         string `mapstructure:"REDIS_ADDRESS"`
	Port         string `mapstructure:"REDIS_PORT"`
	Password     string `mapstructure:"REDIS_PASSWORD"`
	DbNumberAuth string `mapstructure:"REDIS_DB_AUTH"`
}

func NewRedisConfig() (*RedisConfig, error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("REDIS_ADDRESS")
	_ = viper.BindEnv("REDIS_PORT")
	_ = viper.BindEnv("REDIS_PASSWORD")
	_ = viper.BindEnv("REDIS_DB_AUTH")

	var redisConfig RedisConfig
	if err := viper.Unmarshal(&redisConfig); err != nil {
		return nil, err
	}
	return &redisConfig, nil
}
