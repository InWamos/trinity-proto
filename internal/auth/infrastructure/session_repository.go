package infrastructure

import (
	"context"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
)

type SessionRepository interface {
	GetSessionByID(ctx context.Context, id uuid.UUID) (domain.Session, error)
	RevokeSessionByID(ctx context.Context, id uuid.UUID) error
	GetAllSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)
}
