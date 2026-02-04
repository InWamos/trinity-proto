package record

import (
	"github.com/InWamos/trinity-proto/internal/record/domain/telegram/service"
	"go.uber.org/fx"
)

func NewRecordDomainContainer() fx.Option {
	return fx.Module(
		"record_domain",
		fx.Provide(service.NewTelegramModelValidator),
	)
}
