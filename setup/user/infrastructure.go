package user

import (
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
		),
	)
}
