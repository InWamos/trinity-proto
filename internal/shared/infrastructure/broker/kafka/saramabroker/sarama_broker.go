package saramabroker

import (
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
	"github.com/InWamos/trinity-proto/config"
)

// UserSyncProducer We need custom types here for Uber FX.
type UserSyncProducer sarama.SyncProducer
type AuthConsumer sarama.ConsumerGroup

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
	userSyncProducer  UserSyncProducer
	authConsumerGroup AuthConsumer
	logger            *slog.Logger
}

func NewSaramaBroker(config *config.KafkaConfig, logger *slog.Logger) *SaramaBroker {
	brokerLogger := logger.With(
		slog.String("component", "sarama_broker"),
	)
	configKafka := sarama.NewConfig()
	configKafka.Producer.RequiredAcks = sarama.WaitForAll
	configKafka.Producer.Retry.Max = 5
	configKafka.Producer.Return.Successes = true
	configKafka.Consumer.Offsets.Initial = sarama.OffsetOldest

	brokers := unmarshalBrokersSeparatedByComma(config.Brokers)
	userSyncProducer, err := sarama.NewSyncProducer(brokers, configKafka)
	if err != nil {
		brokerLogger.Error("Could not connect to user sync producer", slog.Any("err", err))
		panic(err)
	}

	authConsumerGroup, err := sarama.NewConsumerGroup(brokers, "auth-service-group", configKafka)
	if err != nil {
		brokerLogger.Error("Could not connect to auth consumer group", slog.Any("err", err))
		panic(err)
	}

	return &SaramaBroker{userSyncProducer: userSyncProducer, authConsumerGroup: authConsumerGroup, logger: brokerLogger}
}

func (sb *SaramaBroker) GetSyncUserProducer() UserSyncProducer {
	return sb.userSyncProducer
}

func (sb *SaramaBroker) GetAuthConsumerGroup() AuthConsumer {
	return sb.authConsumerGroup
}

func (sb *SaramaBroker) Close() error {
	if err := sb.userSyncProducer.Close(); err != nil {
		sb.logger.Error("failed to close user sync producer", slog.Any("err", err))
	}
	if err := sb.authConsumerGroup.Close(); err != nil {
		sb.logger.Error("failed to close auth consumer group", slog.Any("err", err))
	}
	return nil
}
