package application

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
)

type GetLatestTelegramRecordsByUserTelegramID struct {
	transactionManagerFactory interfaces.TransactionManagerFactory
	telegramRecordFactory     repository.TelegramRecordRepositoryFactory
	logger                    *slog.Logger
}

type GetLatestTelegramRecordsByUserTelegramIDRequest struct {
}
