package application

import (
	"log/slog"

	domain "github.com/InWamos/trinity-proto/internal/record/domain/telegram"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
)

type GetLatestTelegramRecordsByUserTelegramID struct {
	transactionManagerFactory interfaces.TransactionManagerFactory
	telegramRecordFactory     repository.TelegramRecordRepositoryFactory
	logger                    *slog.Logger
}

type GetLatestTelegramRecordsByUserTelegramIDRequest struct {
	UserTelegramID uint64
}

type GetLatestTelegramRecordsByUserTelegramIDResponse struct {
	TelegramRecords []domain.TelegramRecord
}

func NewGetLatestTelegramRecordsByUserTelegramID(
	transactionManagerFactory interfaces.TransactionManagerFactory,
	telegramRecordFactory repository.TelegramRecordRepositoryFactory,
	logger *slog.Logger,
) *GetLatestTelegramRecordsByUserTelegramID {
	iLogger := logger.With(
		slog.String("module", "record"),
		slog.String("name", "get_latest_tg_records_by_user_telegram_id"),
	)
	return &GetLatestTelegramRecordsByUserTelegramID{
		transactionManagerFactory: transactionManagerFactory,
		telegramRecordFactory:     telegramRecordFactory,
		logger:                    iLogger,
	}
}
