package setup //nolint:dupl //Must be similar

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
)

func runMainServer(server *http.Server, listener *net.Listener, logger *slog.Logger) {
	if err := server.Serve(*listener); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start Main server", slog.Any("err", err))
		panic(err)
	}
}

func getMainServerAndListenConfig(listenAddress string, handler http.Handler) (*http.Server, *net.ListenConfig) {
	srv := &http.Server{
		Addr:              listenAddress,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	listenConfig := &net.ListenConfig{KeepAlive: 3 * time.Minute}
	return srv, listenConfig
}

func getMainServerAndHook(listenAddress string, handler http.Handler, logger *slog.Logger) (*http.Server, fx.Hook) {
	srv, listenConfig := getMainServerAndListenConfig(listenAddress, handler)
	return srv, fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := listenConfig.Listen(ctx, "tcp4", srv.Addr)
			if err != nil {
				return err
			}
			go runMainServer(srv, &ln, logger)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	}
}
