package shared

import (
	"github.com/InWamos/trinity-proto/setup/shared/infrastructure/database"
	"github.com/InWamos/trinity-proto/setup/shared/infrastructure/messagebroker"
	"go.uber.org/fx"
)

func NewSharedModuleContainer() fx.Option {
	return fx.Module(
		"shared_module",
		database.NewSqlxDatabaseContainer(),
		messagebroker.NewSaramaBrokerContainer(),
	)
}
