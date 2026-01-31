package repository

import (
	"context"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
)

type TelegramUserRepository interface {
	GetByTelegramID(ctx context.Context, telegramID uint64) (*domain.TelegramUser, error)
	AddUser(ctx context.Context, user *domain.TelegramUser) error
	DeleteUserByTelegramID(ctx context.Context, telegramID uint64) error
}
