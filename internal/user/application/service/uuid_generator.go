package service

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func (ug *UUIDGenerator) GetUUIDv7() (uuid.UUID, error) {
	return uuid.NewV7()
}
