package repository

import "github.com/InWamos/trinity-proto/internal/shared/interfaces"

type TelegramIdentityRepositoryFactory interface {
	CreateTelegramIdentityRepositoryWithTransaction(tm interfaces.TransactionManager) TelegramIdentityRepository
}
