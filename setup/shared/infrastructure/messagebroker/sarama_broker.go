package messagebroker

import (
	"context"

	"github.com/InWamos/trinity-proto/internal/auth/infrastructure/messaging/saramahandlers"
	"github.com/InWamos/trinity-proto/internal/shared/infrastructure/broker/kafka/saramabroker"
	"go.uber.org/fx"
)

func NewSaramaBrokerContainer() fx.Option {
	return fx.Module(
		"sarama_broker",
		fx.Provide(
			saramabroker.NewSaramaBroker,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, broker *saramabroker.SaramaBroker, authConsumerGroup saramabroker.AuthConsumer, authConsumerHandler *saramahandlers.AuthConsumerGroupHandler) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						authHandlerTopics := []string{"user-events"}
						go authConsumerGroup.Consume( //nolint: errcheck // check for err value is not needed
							context.Background(),
							authHandlerTopics,
							authConsumerHandler,
						)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return broker.Close()
					},
				})
			},
		),
		fx.Provide(
			func(broker *saramabroker.SaramaBroker) saramabroker.UserSyncProducer {
				return broker.GetSyncUserProducer()
			},
			func(broker *saramabroker.SaramaBroker) saramabroker.AuthConsumer {
				return broker.GetAuthConsumerGroup()
			},
		),
	)
}
