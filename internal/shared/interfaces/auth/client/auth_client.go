package client

import (
	"context"
	"errors"
)

var (
	ErrSessionInvalid  = errors.New("session invalid or not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionRevoked  = errors.New("session revoked")
	ErrUnexpectedError = errors.New("unexpected error occurred")
)

type AuthClient interface {
	ValidateSession(ctx context.Context, session string) error
}
