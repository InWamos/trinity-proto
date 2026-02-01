package repository

import (
	"context"
	"errors"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/google/uuid"
)

var (
	ErrIdentityNotFound    = errors.New("Identity not found")
	ErrFailedToAddIdentity = errors.New("Failed to add identity to a database")
)

type TelegramIdentityRepository interface {
	AddIdentity(ctx context.Context, identity *domain.TelegramIdentity) error
	RemoveIdentityByID(ctx context.Context, identityID uuid.UUID) error
	GetIdentityByID(ctx context.Context, identityID uuid.UUID) (*domain.TelegramIdentity, error)
}
