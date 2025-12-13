package user

import (
	sqlxdatabase "github.com/InWamos/trinity-proto/internal/user/infrastructure/database/sqlx_database"
	sqlxrepository "github.com/InWamos/trinity-proto/internal/user/infrastructure/repository/sqlx_repository"
	"go.uber.org/fx"
)

func NewUserInfrastructureContainer() fx.Option {
	return fx.Module(
		"user_infrastructure",
		fx.Provide(
			// Provides Sqlx mapper for repository
			sqlxrepository.NewSqlxMapper,
			// Provides User repository factory
			sqlxrepository.NewSqlxUserRepositoryFactory,
			// Provides GORM sessiona
			sqlxdatabase.NewSQLXTransactionFactory,
			// Provides GormDatabase
			sqlxdatabase.NewSQLXDatabase,
		),
	)
}
