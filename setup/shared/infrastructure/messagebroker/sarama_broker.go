package messagebroker

import (
	"github.com/InWamos/trinity-proto/internal/shared/infrastructure/broker/kafka/saramabroker"
	"go.uber.org/fx"
)

func NewSqlxDatabaseContainer() fx.Option {
	return fx.Module(
		"sarama_broker",
		fx.Provide(
			saramabroker.NewSaramaBroker,
			func(broker *saramabroker.SaramaBroker) saramabroker.UserSyncProducer {
				return broker.GetSyncUserProducer()
			},
			func(broker *saramabroker.SaramaBroker) saramabroker.AuthConsumer {
				return broker.GetAuthConsumerGroup()
			},
		))
}
