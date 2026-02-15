package repositories

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/jmoiron/sqlx"
)

type SQLXTelegramIdentityRepositoryFactory struct {
	logger     *slog.Logger
	sqlxMapper *mappers.SqlxTelegramIdentityMapper
}

func NewSQLXTelegramIdentityRepositoryFactory(
	logger *slog.Logger,
	mapper *mappers.SqlxTelegramIdentityMapper,
) repository.TelegramIdentityRepositoryFactory {
	return &SQLXTelegramIdentityRepositoryFactory{
		logger:     logger,
		sqlxMapper: mapper,
	}
}

func (factory *SQLXTelegramIdentityRepositoryFactory) CreateTelegramIdentityRepositoryWithTransaction(
	tm interfaces.TransactionManager,
) repository.TelegramIdentityRepository {
	tx, ok := tm.GetTransaction().(*sqlx.Tx)
	if !ok {
		factory.logger.Error("invalid transaction type, expected *sqlx.Tx")
		panic("invalid transaction type for sqlx repository")
	}
	return &SQLXTelegramIdentityRepository{
		session:    tx,
		sqlxMapper: factory.sqlxMapper,
		logger:     factory.logger,
	}
}
