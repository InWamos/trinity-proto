package setup

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/rs/cors"
	"go.uber.org/fx"
)

func respondPong(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"hello\": \"world\"}"))
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

func runServer(server *http.Server, listener *net.Listener, logger *slog.Logger) {
	if err := server.Serve(*listener); err != nil {
		logger.Error("Failed to start server", slog.Any("err", err))
		panic(err)
	}
}

func NewHTTPServer(lc fx.Lifecycle, serverConfig *config.ServerConfig, logger *slog.Logger) *http.Server {
	listenAddress := fmt.Sprintf("%s:%d", serverConfig.BindAddress, serverConfig.Port)
	cors := getCORSHeaders(serverConfig.AllowedOrigin)
	masterMux := http.NewServeMux()
	masterMux.HandleFunc("GET /ping", respondPong)

	srv := &http.Server{
		Addr:              listenAddress,
		Handler:           cors.Handler(masterMux),
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
