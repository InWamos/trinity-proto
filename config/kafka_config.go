package config

import "github.com/spf13/viper"

type KafkaConfig struct {
	Brokers          string `mapstructure:"KAFKA_BROKERS"`
	ConsumerGroup    string `mapstructure:"KAFKA_CONSUMER_GROUP"`
	SecurityProtocol string `mapstructure:"KAFKA_SECURITY_PROTOCOL"`
	SASLMechanism    string `mapstructure:"KAFKA_SASL_MECHANISM"`
	SASLUsername     string `mapstructure:"KAFKA_SASL_USERNAME"`
	SASLPassword     string `mapstructure:"KAFKA_SASL_PASSWORD"`
	CompressionType  string `mapstructure:"KAFKA_COMPRESSION_TYPE"`
	RequestTimeoutMs int    `mapstructure:"KAFKA_REQUEST_TIMEOUT_MS"`
}

func NewKafkaConfig() (*KafkaConfig, error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("KAFKA_BROKERS")
	_ = viper.BindEnv("KAFKA_CONSUMER_GROUP")
	_ = viper.BindEnv("KAFKA_SECURITY_PROTOCOL")
	_ = viper.BindEnv("KAFKA_SASL_MECHANISM")
	_ = viper.BindEnv("KAFKA_SASL_USERNAME")
	_ = viper.BindEnv("KAFKA_SASL_PASSWORD")
	_ = viper.BindEnv("KAFKA_COMPRESSION_TYPE")
	_ = viper.BindEnv("KAFKA_REQUEST_TIMEOUT_MS")

	var kafkaConfig KafkaConfig
	if err := viper.Unmarshal(&kafkaConfig); err != nil {
		return nil, err
	}
	return &kafkaConfig, nil
}
