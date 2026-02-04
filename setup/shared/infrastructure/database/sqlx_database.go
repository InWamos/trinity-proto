package database

import (
	sqlxdatabase "github.com/InWamos/trinity-proto/internal/shared/infrastructure/database/sqlx_database"
	"go.uber.org/fx"
)

func NewSqlxDatabaseContainer() fx.Option {
	return fx.Module(
		"sqlx_database",
		fx.Provide(
			sqlxdatabase.NewSQLXDatabase,
			sqlxdatabase.NewSQLXTransactionFactory,
		))
}
