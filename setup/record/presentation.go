package record

import (
	v1 "github.com/InWamos/trinity-proto/internal/record/presentation/v1"
	"github.com/InWamos/trinity-proto/internal/record/presentation/v1/handlers"
	"go.uber.org/fx"
)

func NewRecordPresentationContainer() fx.Option {
	return fx.Module(
		"record_presentation",
		fx.Provide(
			handlers.NewGetLatestTelegramRecordsByTelegramID,
			v1.NewRecordMuxV1,
		),
	)
}
