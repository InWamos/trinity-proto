package saramabroker

import (
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
	"github.com/InWamos/trinity-proto/config"
)

type UserSyncProducer sarama.SyncProducer

func unmarshalBrokersSeparatedByComma(brokers string) []string {
	if brokers == "" {
		return []string{}
	}

	parts := strings.Split(brokers, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

type SaramaBroker struct {
	userSyncProducer UserSyncProducer
	logger           *slog.Logger
}

func NewSaramaBroker(config *config.KafkaConfig, logger *slog.Logger) *SaramaBroker {
	brokerLogger := logger.With(
		slog.String("component", "sarama_broker"),
	)
	configKafka := sarama.NewConfig()
	configKafka.Producer.RequiredAcks = sarama.WaitForAll
	configKafka.Producer.Retry.Max = 5
	configKafka.Producer.Return.Successes = true
	userBrokers := unmarshalBrokersSeparatedByComma(config.Brokers)
	userSyncProducer, err := sarama.NewSyncProducer(userBrokers, configKafka)
	if err != nil {
		brokerLogger.Error("Could not connect to user sync producer", slog.Any("err", err))
		panic(err)
	}
	defer userSyncProducer.Close()

	return &SaramaBroker{userSyncProducer: userSyncProducer, logger: brokerLogger}
}

func (sb *SaramaBroker) GetSyncUserProducer() UserSyncProducer {
	return sb.userSyncProducer
}
