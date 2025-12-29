package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelegramUser struct {
	ID              uuid.UUID `validate:"required,uuid"`
	TelegramID      uint64    `validate:"required,gt=0,lte=300000000000"`
	CurrentUsername string    `validate:"required,min=4,max=32"`
	Bio             string    `validate:"max=200"`
	AddedAt         time.Time `validate:"required"`
	AddedByUser     uuid.UUID `validate:"uuid"`
}
