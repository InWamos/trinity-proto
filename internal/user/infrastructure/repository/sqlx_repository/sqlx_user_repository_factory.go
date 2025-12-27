package sqlxrepository

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/jmoiron/sqlx"
)

type SqlxUserRepositoryFactory struct {
	logger     *slog.Logger
	sqlxMapper *SqlxMapper
}

func NewSqlxUserRepositoryFactory(logger *slog.Logger, mapper *SqlxMapper) repository.UserRepositoryFactory {
	return &SqlxUserRepositoryFactory{
		logger:     logger,
		sqlxMapper: mapper,
	}
}

func (surf *SqlxUserRepositoryFactory) CreateUserRepositoryWithTransaction(
	tm interfaces.TransactionManager,
) repository.UserRepository {
	// Extract the underlying sqlx transaction
	tx, ok := tm.GetTransaction().(*sqlx.Tx)
	if !ok {
		surf.logger.Error("invalid transaction type, expected *sqlx.Tx")
		panic("invalid transaction type for sqlx repository")
	}

	return &SqlxUserRepository{
		session:    tx,
		sqlxMapper: surf.sqlxMapper,
		logger:     surf.logger,
	}
}
