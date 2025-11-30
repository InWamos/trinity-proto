package setup

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"go.uber.org/fx"
)

func getLogger(level string) *slog.Logger {
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		logLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func getConfig() (*config.AppConfig, error) {
	config, err := config.NewAppConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getCORSHeaders(allowedOrigin string) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "DELETE", "PUT", "PATCH", "POST"},
		AllowedHeaders:   []string{"Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(time.Hour.Seconds() * 24),
	})

}

func bootstrapServer(server *gin.Engine, allowOrigin string, trustedProxy string, logger *slog.Logger) {

}

func runServer(server *gin.Engine, bindIP string, port int, logger *slog.Logger) {
	bindAddress := bindIP + ":" + string(rune(port))
	if err := server.Run(bindAddress); err != nil {
		logger.Error("Failed to start server", "error", err, "bind_address", bindAddress)
		panic(err)
	}
}

func NewHTTPServer(lc fx.Lifecycle, serverConfig *config.AppConfig) *http.Server {
	listenAddress := fmt.Sprintf("%s:%d", serverConfig.GinConfig.BindAddress, serverConfig.GinConfig.Port)
	srv := &http.Server{Addr: listenAddress, ReadHeaderTimeout: 5 * time.Second}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lc := net.ListenConfig{}
			ln, err := lc.Listen(ctx, "tcp4", srv.Addr)
			if err != nil {
				return err
			}
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
