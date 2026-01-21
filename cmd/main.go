package main

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/config"
	_ "github.com/InWamos/trinity-proto/docs"
	"github.com/InWamos/trinity-proto/logger"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/InWamos/trinity-proto/setup"
	"github.com/InWamos/trinity-proto/setup/auth"
	"github.com/InWamos/trinity-proto/setup/record"
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
//	@BasePath	/api

//	@securityDefinitions.apikey	SessionCookie
//	@in							cookie
//	@name						session_id
//	@description				Session cookie for authenticated requests. Roles: Admin, User

func main() {
	fx.New(
		fx.Provide(config.NewDatabaseConfig, config.NewLoggingConfig,
			config.NewServerConfig, config.NewRedisConfig),
		fx.Provide(logger.GetLogger),
		fx.Provide(
			middleware.NewGlobalCORSMiddleware,
			middleware.NewTrustedProxyMiddleware,
			middleware.NewLoggingMiddleware,
			middleware.NewAuthenticationMiddleware,
		),
		user.NewUserModuleContainer(),
		auth.NewAuthModuleContainer(),
		record.NewRecordModuleContainer(),
		fx.Provide(setup.NewMainHTTPServer),
		fx.Provide(setup.NewProfilerHTTPServer),
		fx.Provide(setup.NewHTTPServers),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		fx.Invoke(setup.CreateAdminAccountIfNotExists),
		fx.Invoke(func(servers setup.HTTPServers) {}), //nolint:revive //False positive on Fx syntax
	).Run()
}
