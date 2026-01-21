package record

import (
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/repositories"
	"go.uber.org/fx"
)

func NewRecordInfrastructureContainer() fx.Option {
	return fx.Module(
		"record_infrastructure",
		fx.Provide(
			mappers.NewSqlxTelegramRecordMapper,
			repositories.NewSQLXTelegramRecordRepository,
			repositories.NewSQLXTelegramRecordRepositoryFactory,
		),
	)
}
