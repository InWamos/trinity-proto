package main

import (
	"log/slog"
	"net/http"

	"github.com/InWamos/trinity-proto/config"
	"github.com/InWamos/trinity-proto/logger"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/InWamos/trinity-proto/setup"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		fx.Provide(config.NewDatabaseConfig, config.NewLoggingConfig, config.NewServerConfig),
		fx.Provide(logger.GetLogger),
		fx.Provide(middleware.NewGlobalCORSMiddleware, middleware.NewTrustedProxyMiddleware),
		fx.Provide(setup.NewHTTPServer),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
