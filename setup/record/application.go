package record

import (
	application "github.com/InWamos/trinity-proto/internal/record/application/telegram"
	identityApplication "github.com/InWamos/trinity-proto/internal/record/application/telegram/identity"
	"go.uber.org/fx"
)

func NewRecordApplicationContainer() fx.Option {
	return fx.Module(
		"record_application",
		fx.Provide(
			application.NewGetLatestTelegramRecordsByUserTelegramID,
			application.NewAddTelegramUser,
			identityApplication.NewAddTelegramIdentity,
		),
	)
}
