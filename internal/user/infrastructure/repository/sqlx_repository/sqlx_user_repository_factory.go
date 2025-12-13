package sqlxrepository

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/user/infrastructure/repository"
	"github.com/jmoiron/sqlx"
)

type SqlxUserRepositoryFactory struct {
	logger     *slog.Logger
	sqlxMapper *SqlxMapper
}

func NewSqlxUserRepositoryFactory(logger *slog.Logger) repository.UserRepositoryFactory {
	return &SqlxUserRepositoryFactory{logger: logger}
}

func (surf *SqlxUserRepositoryFactory) CreateUserRepository(session any) repository.UserRepository {
	tx, ok := session.(*sqlx.Tx)
	if !ok {
		surf.logger.Error("invalid session type, expected *sqlx.Tx")
		panic("invalid session type for sqlx repository")
	}

	return &SqlxUserRepository{
		session:    tx,
		sqlxMapper: surf.sqlxMapper,
		logger:     surf.logger,
	}
}
