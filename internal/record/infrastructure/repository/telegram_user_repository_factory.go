package repository

import "github.com/InWamos/trinity-proto/internal/shared/interfaces"

type TelegramUserRepositoryFactory interface {
	CreateTelegramUserRepositoryWithTransaction(tm interfaces.TransactionManager) TelegramUserRepository
}
