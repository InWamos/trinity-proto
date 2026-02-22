package messagebroker

import (
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/InWamos/trinity-proto/internal/events/pb"
	"github.com/InWamos/trinity-proto/internal/shared/infrastructure/broker/kafka/saramabroker"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type SaramaUserMessageBroker struct {
	userSyncProducer saramabroker.UserSyncProducer
	logger           *slog.Logger
}

func NewSaramaUserMessageBroker(
	userSyncProducer saramabroker.UserSyncProducer,
	logger *slog.Logger,
) UserMessageBroker {
	sumb := logger.With(slog.String("component", "sarama_user_message_broker"))
	return &SaramaUserMessageBroker{
		userSyncProducer: userSyncProducer,
		logger:           sumb,
	}
}

func (broker *SaramaUserMessageBroker) PostUserRemovedMessage(userID uuid.UUID) error {
	event := &pb.UserRemovedEvent{
		UserId:    userID.String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	serialized, err := proto.Marshal(event)
	if err != nil {
		broker.logger.Error("failed to serialize UserRemovedEvent", slog.Any("err", err))
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: "user-events",
		Key:   sarama.StringEncoder(userID.String()),
		Value: sarama.ByteEncoder(serialized),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte("UserRemoved"),
			},
			{
				Key:   []byte("event_id"),
				Value: []byte(uuid.New().String()),
			},
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().UTC().Format(time.RFC3339)),
			},
		},
	}
	partition, offset, err := broker.userSyncProducer.SendMessage(message)
	if err != nil {
		broker.logger.Error("failed to send UserRemovedEvent to Kafka",
			slog.String("user_id", userID.String()),
			slog.Any("err", err))
		return err
	}

	broker.logger.Debug("UserRemovedEvent sent to Kafka",
		slog.String("user_id", userID.String()),
		slog.Int64("partition", int64(partition)),
		slog.Int64("offset", offset))

	return nil
}
