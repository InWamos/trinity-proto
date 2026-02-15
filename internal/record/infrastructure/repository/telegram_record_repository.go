package repository

import (
	"context"
	"errors"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
)

var (
	ErrTelegramRecordCreationFailed = errors.New("failed to save telegram record")
	ErrDatabaseFailed               = errors.New("database request has failed")
)

type TelegramRecordRepository interface {
	GetLatestTelegramRecordsByUserTelegramID(
		ctx context.Context,
		userTelegramID uint64,
	) (*[]domain.TelegramRecord, error)
	CreateTelegramRecord(ctx context.Context, telegramRecord domain.TelegramRecord) error
	CreateTelegramRecords(ctx context.Context, telegramRecords []domain.TelegramRecord) error
}
