package repositories

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/jmoiron/sqlx"
)

type SQLXTelegramUserRepositoryFactory struct {
	logger     *slog.Logger
	sqlxMapper *mappers.SqlxTelegramUserMapper
}

func NewSQLXTelegramUserRepositoryFactory(
	logger *slog.Logger,
	mapper *mappers.SqlxTelegramUserMapper,
) repository.TelegramUserRepositoryFactory {
	return &SQLXTelegramUserRepositoryFactory{
		logger:     logger,
		sqlxMapper: mapper,
	}
}

// TODO: Implement User Telegram identity business logic
func (factory *SQLXTelegramUserRepositoryFactory) CreateTelegramUserRepositoryWithTransaction(
	tm interfaces.TransactionManager,
) repository.TelegramUserRepository {
	tx, ok := tm.GetTransaction().(*sqlx.Tx)
	if !ok {
		factory.logger.Error("invalid transaction type, expected *sqlx.Tx")
		panic("invalid transaction type for sqlx repository")
	}
	return &SQLXTelegramUserRepository{
		session:    tx,
		sqlxMapper: factory.sqlxMapper,
		logger:     factory.logger,
	}
}
