package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelegramRecord struct {
	ID          uuid.UUID `validate:"required,uuid"`
	Author      uuid.UUID `validate:"required,uuid"`
	MessageText string    `validate:"required,min=1,max=4096"`
	Attachments []uuid.UUID
	PostedAt    time.Time `validate:"required"`
	AddedByUser uuid.UUID `validate:"required,uuid"`
}
