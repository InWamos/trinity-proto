package setup

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/InWamos/trinity-proto/config"
	authV1Mux "github.com/InWamos/trinity-proto/internal/auth/presentation/v1"
	recordV1Mux "github.com/InWamos/trinity-proto/internal/record/presentation/v1"
	"github.com/InWamos/trinity-proto/internal/user/application"
	userV1Mux "github.com/InWamos/trinity-proto/internal/user/presentation/v1"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
)

var ErrProfilerEnabledInNonDevelopment = errors.New("you MUST disable profiler in non-dev environment")

type MainHTTPServer struct {
	*http.Server
}

type ProfilerHTTPServer struct {
	*http.Server
}

type HTTPServers struct {
	Main     MainHTTPServer
	Profiler ProfilerHTTPServer
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

func NewMainHTTPServer(
	lc fx.Lifecycle,
	serverConfig *config.ServerConfig,
	loggingMiddleware *middleware.LoggingMiddleware,
	corsMiddleware *middleware.GlobalCORSMiddleware,
	trustedProxyMiddleware *middleware.TrustedProxyMiddleware,
	authMiddleware *middleware.AuthenticationMiddleware,
	userMuxV1 *userV1Mux.UserMuxV1,
	authMuxV1 *authV1Mux.AuthMuxV1,
	sessionManagementMuxV1 *authV1Mux.SessionManagementMuxV1,
	recordMuxV1 *recordV1Mux.RecordMuxV1,
	logger *slog.Logger,
) MainHTTPServer {
	listenAddress := fmt.Sprintf("%s:%d", serverConfig.BindAddress, serverConfig.Port)
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
	chiRouter.Mount("/api/v1/record", authMiddleware.Handler(recordMuxV1.GetMux()))
	chiRouter.Mount("/api/v1/auth", authMuxV1.GetMux())
	chiRouter.Mount("/api/v1/session", authMiddleware.Handler(sessionManagementMuxV1.GetMux()))
	chiRouter.Mount("/swagger", httpSwagger.WrapHandler)

	srv, hook := getMainServerAndHook(listenAddress, chiRouter, logger)
	// Add Uber Fx OnStart and OnStop hooks.
	// Help to organize server startup and shutdown by clear, net/http native startup and shutdown operation.
	lc.Append(hook)
	return MainHTTPServer{srv}
}

// NewProfilerHTTPServer adds a lifecycle for Pprof server on port 6060.
// Only creates and registers the profiler in development environment.
// Returns nil in production.
func NewProfilerHTTPServer(
	lc fx.Lifecycle,
	serverConfig *config.ServerConfig,
	logger *slog.Logger,
) ProfilerHTTPServer {
	if serverConfig.Environment != "DEVELOPMENT" {
		logger.Info("Profiler server disabled in non-development environment")
		return ProfilerHTTPServer{nil}
	}
	profilerListenAddress := fmt.Sprintf("%s:%d", serverConfig.BindAddress, 6060)
	srv, hook := getProfilerServerAndHook(profilerListenAddress, nil, logger)
	lc.Append(hook)
	return ProfilerHTTPServer{srv}
}

// NewHTTPServers provides both main and profiler HTTP servers.
func NewHTTPServers(main MainHTTPServer, profiler ProfilerHTTPServer) HTTPServers {
	return HTTPServers{
		Main:     main,
		Profiler: profiler,
	}
}
