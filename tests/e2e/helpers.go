package e2e

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/InWamos/trinity-proto/logger"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/InWamos/trinity-proto/setup"
	"github.com/InWamos/trinity-proto/setup/user"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
)

// VerboseFxLogger logs every fx event including constructor calls
type VerboseFxLogger struct {
	logger *slog.Logger
}

func (l *VerboseFxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Debug("fx lifecycle: OnStart hook executing",
			slog.String("caller", e.FunctionName),
			slog.String("callee", e.CallerName))
	case *fxevent.OnStartExecuted:
		l.logger.Debug("fx lifecycle: OnStart hook executed",
			slog.String("caller", e.FunctionName),
			slog.String("callee", e.CallerName),
			slog.Duration("runtime", e.Runtime))
	case *fxevent.OnStopExecuting:
		l.logger.Debug("fx lifecycle: OnStop hook executing",
			slog.String("caller", e.FunctionName),
			slog.String("callee", e.CallerName))
	case *fxevent.OnStopExecuted:
		l.logger.Debug("fx lifecycle: OnStop hook executed",
			slog.String("caller", e.FunctionName),
			slog.String("callee", e.CallerName),
			slog.Duration("runtime", e.Runtime))
	case *fxevent.Supplied:
		l.logger.Debug("fx: constructor supplied",
			slog.String("type", e.TypeName))
	case *fxevent.Provided:
		for _, output := range e.OutputTypeNames {
			l.logger.Debug("fx: object created",
				slog.String("constructor", e.ConstructorName),
				slog.String("type", output),
				slog.Bool("private", e.Private))
		}
	case *fxevent.Invoked:
		l.logger.Debug("fx: function invoked",
			slog.String("function", e.FunctionName))
	case *fxevent.Invoking:
		l.logger.Debug("fx: invoking function",
			slog.String("function", e.FunctionName))
	case *fxevent.Started:
		l.logger.Debug("fx: application started",
			slog.Any("error", e.Err))
	case *fxevent.Stopped:
		l.logger.Debug("fx: application stopped",
			slog.Any("error", e.Err))
	}
}

// StartTestServer starts the fx application for testing and returns a cleanup function
func StartTestServer(t *testing.T) (baseURL string, cleanup func()) {
	t.Helper()

	app := fxtest.New(t,
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
			return &VerboseFxLogger{logger: logger}
		}),
		fx.Invoke(func(*http.Server) {}),
	)

	app.RequireStart()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	cleanup = func() {
		t.Helper()
		t.Log("Stopping test server...")

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Stop the fx app, which will call srv.Shutdown(ctx)
		if err := app.Stop(ctx); err != nil {
			t.Logf("Warning: error stopping app: %v", err)
		}

		// Wait a bit to ensure the port is fully released
		time.Sleep(200 * time.Millisecond)
		t.Log("Test server stopped")
	}

	baseURL = "http://127.0.0.1:18080"
	return baseURL, cleanup
}
