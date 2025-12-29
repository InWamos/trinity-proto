package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelegramRecord struct {
	ID          uuid.UUID
	Author      uuid.UUID
	MessageText string
	Attachments []uuid.UUID
	PostedAt    time.Time
	AddedByUser uuid.UUID
}
