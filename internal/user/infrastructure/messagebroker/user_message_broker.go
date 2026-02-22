package messagebroker

import (
	"github.com/google/uuid"
)

type UserMessageBroker interface {
	PostUserRemovedMessage(userID uuid.UUID) error
}
