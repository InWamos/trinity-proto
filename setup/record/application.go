package record

import (
	application "github.com/InWamos/trinity-proto/internal/record/application/telegram"
	identityApplication "github.com/InWamos/trinity-proto/internal/record/application/telegram/identity"
	"github.com/InWamos/trinity-proto/internal/record/application/telegram/record"
	"go.uber.org/fx"
)

func NewRecordApplicationContainer() fx.Option {
	return fx.Module(
		"record_application",
		fx.Provide(
			application.NewGetLatestTelegramRecordsByUserTelegramID,
			application.NewAddTelegramUser,
			record.NewAddTelegramRecord,
			identityApplication.NewAddTelegramIdentity,
		),
	)
}
