package client

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSessionInvalid  = errors.New("session invalid or not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionRevoked  = errors.New("session revoked")
	ErrUnexpectedError = errors.New("unexpected error occurred")
)

type UserRole string

const (
	Admin UserRole = "admin"
	User  UserRole = "user"
)

type UserIdentity struct {
	UserID   uuid.UUID
	UserRole UserRole
}

type AuthClient interface {
	ValidateSession(ctx context.Context, token string) (UserIdentity, error)
}
