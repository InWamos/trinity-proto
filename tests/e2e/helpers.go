package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/InWamos/trinity-proto/logger"
	"github.com/InWamos/trinity-proto/middleware"
	"github.com/InWamos/trinity-proto/setup"
	"github.com/InWamos/trinity-proto/setup/auth"
	"github.com/InWamos/trinity-proto/setup/record"
	"github.com/InWamos/trinity-proto/setup/shared"
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
		fx.Provide(config.NewDatabaseConfig, config.NewLoggingConfig, config.NewServerConfig, config.NewRedisConfig),
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
		shared.NewSharedModuleContainer(),
		fx.Provide(setup.NewMainHTTPServer),
		fx.Provide(setup.NewProfilerHTTPServer),
		fx.Provide(setup.NewHTTPServers),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &VerboseFxLogger{logger: logger}
		}),
		fx.Invoke(func(servers setup.HTTPServers) {}), //nolint:revive //False positive on Fx syntax
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

// LoginResponse represents the response from the login endpoint
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// LoginUser logs in a user and returns an authorization token
func LoginUser(t *testing.T, baseURL, username, password string) string {
	t.Helper()

	loginBody := map[string]string{
		"username": username,
		"password": password,
	}
	body, err := json.Marshal(loginBody)
	if err != nil {
		t.Fatalf("failed to marshal login request: %v", err)
	}

	resp, err := http.Post(
		baseURL+"/api/v1/auth/login",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read login response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		t.Fatalf("failed to unmarshal login response: %v", err)
	}

	if loginResp.Token == "" {
		t.Fatal("login response missing token")
	}

	return loginResp.Token
}

// MakeAuthorizedRequest makes an HTTP request with authorization header
func MakeAuthorizedRequest(t *testing.T, method, url, token string, body interface{}) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	return resp
}
