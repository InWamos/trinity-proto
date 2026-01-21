package repository

import "github.com/InWamos/trinity-proto/internal/shared/interfaces"

type TelegramRecordRepositoryFactory interface {
	CreateTelegramRecordRepositoryWithTransaction(tm interfaces.TransactionManager) TelegramRecordRepository
}
