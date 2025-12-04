package main

import (
	"net/http"

	"github.com/InWamos/trinity-proto/config"
	"github.com/InWamos/trinity-proto/logger"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/InWamos/trinity-proto/setup"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.NewDatabaseConfig, config.NewLoggingConfig, config.NewServerConfig),
		fx.Provide(logger.GetLogger),
		fx.Provide(middleware.NewGlobalCORSMiddleware, middleware.NewTrustedProxyMiddleware),
		fx.Provide(setup.NewHTTPServer),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
