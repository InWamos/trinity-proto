package auth

import (
	authclient "github.com/InWamos/trinity-proto/internal/auth/presentation/client"
	userclient "github.com/InWamos/trinity-proto/internal/user/presentation/client"
	"go.uber.org/fx"
)

func NewAuthPresentationContainer() fx.Option {
	return fx.Module(
		"auth_presentation",
		fx.Provide(
			// Provides user client for credential verification
			userclient.NewUserClient,
			// Provides auth client for session validation
			authclient.NewAuthClient,
		),
	)
}
