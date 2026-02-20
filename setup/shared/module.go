package shared

import (
	"github.com/InWamos/trinity-proto/setup/shared/infrastructure/database"
	"go.uber.org/fx"
)

func NewSharedModuleContainer() fx.Option {
	return fx.Module(
		"shared_module",
		database.NewSqlxDatabaseContainer(),
	)
}
