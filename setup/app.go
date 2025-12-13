package setup

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/InWamos/trinity-proto/config"
	userV1Mux "github.com/InWamos/trinity-proto/internal/user/presentation/v1"
	"github.com/InWamos/trinity-proto/middleware"
	"go.uber.org/fx"
)

func runServer(server *http.Server, listener *net.Listener, logger *slog.Logger) {
	if err := server.Serve(*listener); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", slog.Any("err", err))
		panic(err)
	}
}

func NewHTTPServer(
	lc fx.Lifecycle,
	serverConfig *config.ServerConfig,
	loggingMiddleware *middleware.LoggingMiddleware,
	corsMiddleware *middleware.GlobalCORSMiddleware,
	trustedProxyMiddleware *middleware.TrustedProxyMiddleware,
	userMuxV1 *userV1Mux.UserMuxV1,
	logger *slog.Logger,
) *http.Server {
	listenAddress := fmt.Sprintf("%s:%d", serverConfig.BindAddress, serverConfig.Port)
	masterMux := http.NewServeMux()
	masterMux.Handle("/api/v1/users", userMuxV1.GetMux())
	masterHandler := loggingMiddleware.Handler(corsMiddleware.Handler(trustedProxyMiddleware.Handler(masterMux)))

	srv := &http.Server{
		Addr:              listenAddress,
		Handler:           masterHandler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	listenConfig := &net.ListenConfig{KeepAlive: 3 * time.Minute}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := listenConfig.Listen(ctx, "tcp4", srv.Addr)
			if err != nil {
				return err
			}
			go runServer(srv, &ln, logger)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
