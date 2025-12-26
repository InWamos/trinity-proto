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
	authV1Mux "github.com/InWamos/trinity-proto/internal/auth/presentation/v1"
	"github.com/InWamos/trinity-proto/internal/user/application"
	userV1Mux "github.com/InWamos/trinity-proto/internal/user/presentation/v1"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
)

func runServer(server *http.Server, listener *net.Listener, logger *slog.Logger) {
	if err := server.Serve(*listener); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start server", slog.Any("err", err))
		panic(err)
	}
}
func runProfiler(listenAddress string, logger *slog.Logger) {
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		logger.Error("Profiler server error", slog.Any("err", err))
		panic(err)
	}
}

func CreateAdminAccountIfNotExists(
	interactor *application.CreateRandomAdminUser,
	logger *slog.Logger,
) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Info("Checking for admin account...")
	if err := interactor.Execute(ctx); err != nil {
		logger.Error("Failed to create admin account", slog.Any("error", err))
		panic(err)
	}
	logger.Info("Admin account check complete")
}

func NewHTTPServer(
	lc fx.Lifecycle,
	serverConfig *config.ServerConfig,
	loggingMiddleware *middleware.LoggingMiddleware,
	corsMiddleware *middleware.GlobalCORSMiddleware,
	trustedProxyMiddleware *middleware.TrustedProxyMiddleware,
	authMiddleware *middleware.AuthenticationMiddleware,
	userMuxV1 *userV1Mux.UserMuxV1,
	authMuxV1 *authV1Mux.AuthMuxV1,
	logger *slog.Logger,
) *http.Server {
	if serverConfig.Environment == "DEVELOPMENT" {
		listenAddress := fmt.Sprintf("%s:%d", serverConfig.BindAddress, 6060)

		go runProfiler(listenAddress, logger)
	}
	listenAddress := fmt.Sprintf("%s:%d", serverConfig.BindAddress, serverConfig.Port)
	// masterMux := http.NewServeMux()
	chiRouter := chi.NewRouter()
	// Logging
	chiRouter.Use(loggingMiddleware.Handler)
	// Only allow json content type
	chiRouter.Use(chiMiddleware.AllowContentType("application/json"))
	// Send nocache header to reverse proxy
	chiRouter.Use(chiMiddleware.NoCache)
	// Real IP spoofing protection
	chiRouter.Use(trustedProxyMiddleware.Handler)
	// Fixes typos in /
	chiRouter.Use(chiMiddleware.CleanPath)
	// Recover from panic, log it and return 500
	chiRouter.Use(chiMiddleware.Recoverer)
	// Trim slash from the end of uri
	chiRouter.Use(chiMiddleware.RedirectSlashes)
	// CORS
	chiRouter.Use(corsMiddleware.Handler)
	chiRouter.Mount("/api/v1/users", authMiddleware.Handler(userMuxV1.GetMux()))
	chiRouter.Mount("/api/v1/auth", authMuxV1.GetMux())
	chiRouter.Mount("/swagger/", httpSwagger.WrapHandler)
	// /api/v1/users set of handlers - protected with auth middleware
	// masterMux.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", authMiddleware.Handler(userMuxV1.GetMux())))
	// /api/v1/auth set of handlers
	// masterMux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authMuxV1.GetMux()))

	// Swagger documentation
	// masterMux.Handle("/swagger/", httpSwagger.WrapHandler)

	// masterHandler := loggingMiddleware.Handler(
	// 	corsMiddleware.Handler(trustedProxyMiddleware.Handler(masterMux)),
	// )

	srv := &http.Server{
		Addr:              listenAddress,
		Handler:           chiRouter,
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
