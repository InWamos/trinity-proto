package infrastructure

import (
	"context"
	"errors"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInternal        = errors.New("internal error")
)

type SessionRepository interface {
	GetSessionByToken(ctx context.Context, token string) (domain.Session, error)
	RevokeSessionByToken(ctx context.Context, token string) error
	RevokeAllSessionsByUserID(ctx context.Context, userID uuid.UUID) error
	GetAllSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)
	CreateSession(ctx context.Context, session domain.Session) error
}
