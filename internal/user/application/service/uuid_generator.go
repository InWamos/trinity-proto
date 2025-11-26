package service

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrFailedToGenerateUUID = errors.New("failed to generate UUID")
)

type UUIDGenerator struct {
}

func (ug *UUIDGenerator) GetUUIDv7() (uuid.UUID, error) {
	result, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, ErrFailedToGenerateUUID
	}
	return result, nil
}
