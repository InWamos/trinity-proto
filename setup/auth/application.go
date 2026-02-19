package auth

import (
	"github.com/InWamos/trinity-proto/internal/auth/application"
	"go.uber.org/fx"
)

func NewAuthApplicationContainer() fx.Option {
	return fx.Module(
		"auth_application",
		fx.Provide(
			// Provides AddSession interactor
			application.NewAddSession,
			// Provides VerifySession interactor
			application.NewVerifySession,
			application.NewRemoveSession,
		),
	)
}
