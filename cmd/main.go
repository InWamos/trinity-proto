package main

import (
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/config"
	"github.com/InWamos/trinity-proto/logger"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/InWamos/trinity-proto/setup"
	"github.com/InWamos/trinity-proto/setup/user"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

//	@title			Trinity API
//	@version		1.0
//	@description	Trinity API schema.

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	MIT

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	SessionCookie
//	@in							cookie
//	@name						session_id
//	@description				Session cookie for authenticated requests. Roles: Admin, User

func main() {
	fx.New(
		fx.Provide(config.NewDatabaseConfig, config.NewLoggingConfig, config.NewServerConfig),
		fx.Provide(logger.GetLogger),
		fx.Provide(
			middleware.NewGlobalCORSMiddleware,
			middleware.NewTrustedProxyMiddleware,
			middleware.NewLoggingMiddleware,
		),
		user.NewUserModuleContainer(),
		fx.Provide(setup.NewHTTPServer),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
