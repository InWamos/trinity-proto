package repository

import (
	"context"
	"errors"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
)

var ErrTelegramRecordCreationFailed = errors.New("failed to save telegram record")

type TelegramRecordRepository interface {
	GetLatestTelegramRecordsByUserTelegramID(
		ctx context.Context,
		userTelegramID uint64,
	) ([]domain.TelegramRecord, error)
	CreateTelegramRecord(ctx context.Context, telegramRecord domain.TelegramRecord) error
	CreateTelegramRecords(ctx context.Context, telegramRecords []domain.TelegramRecord) error
}

type TelegramRecordRepositoryFactory interface {
	CreateTelegramRecordRepositoryWithTransaction(tm interfaces.TransactionManager) TelegramRecordRepository
}
