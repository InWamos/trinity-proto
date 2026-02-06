package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrIdentityNotFound      = errors.New("identity not found")
	ErrIdentityAlreadyExists = errors.New("identity already exists")
)

type TelegramIdentity struct {
	ID          uuid.UUID `validate:"required,uuid"`
	UserID      uuid.UUID `validate:"required,uuid"`
	Username    string    `validate:"required,min=4,max=32"`
	FirstName   string    `validate:"required,min=1,max=64"`
	LastName    string    `validate:"min=0,max=64"`
	Bio         string    `validate:"max=140"`
	PhoneNumber string    `validate:"e164"`
	AddedAt     time.Time `validate:"required"`
	AddedByUser uuid.UUID `validate:"uuid"`
}
