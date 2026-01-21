package repositories

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/jmoiron/sqlx"
)

type SQLXTelegramRecordRepositoryFactory struct {
	logger     *slog.Logger
	sqlxMapper *mappers.SqlxTelegramRecordMapper
}

func NewSQLXTelegramRecordRepositoryFactory(
	logger *slog.Logger,
	mapper *mappers.SqlxTelegramRecordMapper,
) repository.TelegramRecordRepositoryFactory {
	return &SQLXTelegramRecordRepositoryFactory{
		logger:     logger,
		sqlxMapper: mapper,
	}
}

func (factory *SQLXTelegramRecordRepositoryFactory) CreateTelegramRecordRepositoryWithTransaction(
	tm interfaces.TransactionManager,
) repository.TelegramRecordRepository {
	tx, ok := tm.GetTransaction().(*sqlx.Tx)
	if !ok {
		factory.logger.Error("invalid transaction type, expected *sqlx.Tx")
		panic("invalid transaction type for sqlx repository")
	}
	return &SQLXTelegramRecordRepository{
		session:    tx,
		sqlxMapper: factory.sqlxMapper,
		logger:     factory.logger,
	}
}
