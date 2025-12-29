package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelegramUser struct {
	ID              uuid.UUID
	TelegramID      uint64
	CurrentUsername string
	Bio             string
	AddedAt         time.Time
	AddedByUser     uuid.UUID
}
