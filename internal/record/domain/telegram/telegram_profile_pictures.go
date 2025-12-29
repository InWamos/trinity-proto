package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelegramProfilePictures struct {
	ID             uuid.UUID
	UserTelegramID uint64
	StorageKey     string
	PostedAt       time.Time
	MimeType       string
	AddedAt        time.Time
	AddedByUser    uuid.UUID
}
