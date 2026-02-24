package saramahandlers

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/InWamos/trinity-proto/internal/auth/infrastructure/messaging/messageprocessors"
	"github.com/InWamos/trinity-proto/internal/events/pb"
)

type AuthConsumerGroupHandler struct {
	removeUserProcessor *messageprocessors.RemoveUserProcessor
	logger              *slog.Logger
}

func NewAuthConsumerGroupHandler(
	removeUserProcessor *messageprocessors.RemoveUserProcessor,
	logger *slog.Logger,
) *AuthConsumerGroupHandler {
	return &AuthConsumerGroupHandler{
		removeUserProcessor: removeUserProcessor,
		logger:              logger.With(slog.String("component", "auth_consumer_group_handler")),
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *AuthConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	h.logger.Info("Auth consumer group handler session started")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (h *AuthConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.logger.Info("Auth consumer group handler session cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (h *AuthConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				h.logger.Warn("Message channel has been closed")
				return nil
			}

			if err := h.processMessage(session, message); err != nil {
				// Don't mark message on error - let it retry
				continue
			}

			// Mark as consumed only after successful processing
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance.
		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage handles a single Kafka message.
func (h *AuthConsumerGroupHandler) processMessage(
	session sarama.ConsumerGroupSession,
	message *sarama.ConsumerMessage,
) error {
	eventType := h.extractEventType(message.Headers)
	h.logger.DebugContext(session.Context(), "Processing message",
		slog.String("topic", message.Topic),
		slog.String("event_type", eventType))

	// Route to appropriate event processor based on event type
	switch eventType {
	case "UserRemoved":
		event, err := h.deserializeUserRemovedEvent(session.Context(), message.Value)
		if err != nil {
			return err
		}
		return h.removeUserProcessor.Process(session.Context(), event)

	default:
		h.logger.WarnContext(session.Context(), "Unknown event type",
			slog.String("event_type", eventType))
		return nil
	}
}

// extractEventType extracts the event_type from message headers.
func (h *AuthConsumerGroupHandler) extractEventType(headers []*sarama.RecordHeader) string {
	for _, header := range headers {
		if string(header.Key) == "event_type" {
			return string(header.Value)
		}
	}
	return ""
}

// deserializeUserRemovedEvent deserializes protobuf UserRemovedEvent from message bytes.
func (h *AuthConsumerGroupHandler) deserializeUserRemovedEvent(
	ctx context.Context,
	data []byte,
) (*pb.UserRemovedEvent, error) {
	event := &pb.UserRemovedEvent{}
	if err := proto.Unmarshal(data, event); err != nil {
		h.logger.ErrorContext(ctx, "Failed to deserialize UserRemovedEvent",
			slog.Any("err", err))
		return nil, err
	}
	return event, nil
}
