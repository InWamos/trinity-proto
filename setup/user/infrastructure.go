package user

import (
	sqlxdatabase "github.com/InWamos/trinity-proto/internal/shared/infrastructure/database/sqlx_database"
	sqlxrepository "github.com/InWamos/trinity-proto/internal/user/infrastructure/repository/sqlx_repository"
	"go.uber.org/fx"
)

func NewUserInfrastructureContainer() fx.Option {
	return fx.Module(
		"user_infrastructure",
		fx.Provide(
			// Provides Sqlx mapper for repository
			sqlxrepository.NewSqlxUserMapper,
			// Provides User repository factory
			sqlxrepository.NewSqlxUserRepositoryFactory,
			// Provides SQLx session
			sqlxdatabase.NewSQLXTransactionFactory,
			// Provides SQLx Database
			sqlxdatabase.NewSQLXDatabase,
		),
	)
}
