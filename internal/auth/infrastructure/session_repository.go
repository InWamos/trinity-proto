package infrastructure

import (
	"context"
	"errors"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type SessionRepository interface {
	GetSessionByID(ctx context.Context, id uuid.UUID) (domain.Session, error)
	RevokeSessionByID(ctx context.Context, id uuid.UUID) error
	GetAllSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)
}
