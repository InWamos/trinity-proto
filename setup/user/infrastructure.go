package user

import (
	database "github.com/InWamos/trinity-proto/internal/user/infrastructure/database/gorm_database"
	gormrepository "github.com/InWamos/trinity-proto/internal/user/infrastructure/repository/gorm_repository"
	"go.uber.org/fx"
)

func NewUserInfrastructureContainer() fx.Option {
	return fx.Module(
		"user_infrastructure",
		fx.Provide(
			// Provides Gorm mapper for repository
			gormrepository.NewGormMapper,
			// Provides User repository
			gormrepository.NewGormUserRepository,
			// Provides GORM session
			database.NewGormTransactionFactory,
			// Provides GormDatabase
			database.NewGormDatabase,
		),
	)
}
