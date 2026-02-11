package record

import (
	"github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/mappers"
	SqlxTelegramIdentityRepositories "github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/repositories/telegram_identity"
	SqlxTelegramRecordRepositories "github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/repositories/telegram_record"
	SqlxTelegramUserRepositories "github.com/InWamos/trinity-proto/internal/record/infrastructure/repository/sqlx/repositories/telegram_user"
	"go.uber.org/fx"
)

func NewRecordInfrastructureContainer() fx.Option {
	return fx.Module(
		"record_infrastructure",
		fx.Provide(
			mappers.NewSqlxTelegramRecordMapper,
			mappers.NewSqlxTelegramUserMapper,
			SqlxTelegramRecordRepositories.NewSQLXTelegramRecordRepository,
			SqlxTelegramRecordRepositories.NewSQLXTelegramRecordRepositoryFactory,
			SqlxTelegramUserRepositories.NewSQLXTelegramUserRepository,
			SqlxTelegramUserRepositories.NewSQLXTelegramUserRepositoryFactory,
			SqlxTelegramIdentityRepositories.NewSQLXTelegramIdentityRepository,
			SqlxTelegramIdentityRepositories.NewSQLXTelegramIdentityRepositoryFactory,
		),
	)
}
