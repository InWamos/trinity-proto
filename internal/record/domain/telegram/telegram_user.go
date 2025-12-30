package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// TelegramUser Validation domain rules according to https://limits.tginfo.me/en.
type TelegramUser struct {
	ID              uuid.UUID `validate:"required,uuid"`
	TelegramID      uint64    `validate:"required,gt=0,lte=300000000000"`
	CurrentUsername string    `validate:"required,min=4,max=32"`
	FirstName       string    `validate:"min=1,max=64"`
	LastName        string    `validate:"min=0,max=64"`
	Bio             string    `validate:"max=140"`
	AddedAt         time.Time `validate:"required"`
	AddedByUser     uuid.UUID `validate:"uuid"`
}
