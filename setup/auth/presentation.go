package auth

import (
	authclient "github.com/InWamos/trinity-proto/internal/auth/presentation/client"
	authv1mux "github.com/InWamos/trinity-proto/internal/auth/presentation/v1"
	"github.com/InWamos/trinity-proto/internal/auth/presentation/v1/handlers"
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
			// Provides login handler
			handlers.NewLoginHandler,
			// Provides auth multiplexer with routes
			authv1mux.NewAuthMuxV1,
		),
	)
}
