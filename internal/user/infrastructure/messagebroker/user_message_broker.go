package messagebroker

import (
	"context"

	"github.com/google/uuid"
)

type UserMessageBroker interface {
	PostUserRemovedMessage(ctx context.Context, userID uuid.UUID) error
}
