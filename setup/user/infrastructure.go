package user

import (
	"github.com/InWamos/trinity-proto/internal/user/infrastructure/database"
	gormrepository "github.com/InWamos/trinity-proto/internal/user/infrastructure/repository/gorm_repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewUserInfrastructureContainer() fx.Option {
	return fx.Module(
		"user_infrastructure",
		fx.Provide(
			// Provides Gorm mapper for repository
			gormrepository.NewGormMapper,
			// Provides User repository
			gormrepository.NewGormUserRepository,
			// Provides Transaction manager
			database.NewGormTransactionManager,
			// Provides GORM session
			database.NewGormTransaction,
			// Provides GormDatabase
			database.NewGormDatabase,
			// Provides *gorm.DB
			NewGormDB,
		),
	)
}

func NewGormDB(gormDB *database.GormDatabase) *gorm.DB {
	return gormDB.GetEngine()
}
