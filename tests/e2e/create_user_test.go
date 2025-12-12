package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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

var (
	testContainers *TestContainers
	serverBaseURL  string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup containers
	var err error
	testContainers, err = SetupContainers(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup containers: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables for the test
	for k, v := range testContainers.GetDatabaseConfig() {
		os.Setenv(k, v)
	}
	for k, v := range testContainers.GetRedisConfig() {
		os.Setenv(k, v)
	}

	// Server config
	os.Setenv("SERVER_ADDRESS", "127.0.0.1")
	os.Setenv("SERVER_PORT", "18080") // Use different port for tests
	os.Setenv("SERVER_TRUSTED_PROXY", "127.0.0.1")
	os.Setenv("SERVER_ALLOWED_ORIGIN", "*")
	os.Setenv("LOGGING_LEVEL", "debug")

	serverBaseURL = "http://127.0.0.1:18080"

	// Run tests
	code := m.Run()

	// Teardown
	if err := testContainers.Teardown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to teardown containers: %v\n", err)
	}

	os.Exit(code)
}

// startTestServer starts the fx application for testing
func startTestServer(t *testing.T) *fxtest.App {
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
			return &fxevent.SlogLogger{Logger: logger}
		}),
		fx.Invoke(func(*http.Server) {}),
	)

	app.RequireStart()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	return app
}

func TestCreateUser_Success(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Prepare request
	reqBody := map[string]string{
		"username":     "testuser",
		"display_name": "Test User",
		"password":     "password123",
		"role":         "user",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusCreated, resp.StatusCode, string(respBody))
	}

	// Verify response message
	var response map[string]string
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "The user has been created. you can login now"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, response["message"])
	}
}

func TestCreateUser_InvalidUsername_TooShort(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Prepare request with invalid username (too short)
	reqBody := map[string]string{
		"username":     "a", // min is 2
		"display_name": "Test User",
		"password":     "password123",
		"role":         "user",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusBadRequest, resp.StatusCode, string(respBody))
	}
}

func TestCreateUser_InvalidPassword_TooShort(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Prepare request with invalid password (too short)
	reqBody := map[string]string{
		"username":     "testuser2",
		"display_name": "Test User",
		"password":     "short", // min is 8
		"role":         "user",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusBadRequest, resp.StatusCode, string(respBody))
	}
}

func TestCreateUser_InvalidRole(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Prepare request with invalid role
	reqBody := map[string]string{
		"username":     "testuser3",
		"display_name": "Test User",
		"password":     "password123",
		"role":         "superadmin", // only "user" or "admin" allowed
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusBadRequest, resp.StatusCode, string(respBody))
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Prepare request with missing required fields
	reqBody := map[string]string{
		"username": "testuser4",
		// missing display_name, password, role
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusBadRequest, resp.StatusCode, string(respBody))
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Send invalid JSON
	body := []byte(`{"username": "test", invalid json}`)

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusBadRequest, resp.StatusCode, string(respBody))
	}
}

func TestCreateUser_AdminRole(t *testing.T) {
	app := startTestServer(t)
	defer app.RequireStop()

	// Prepare request with admin role
	reqBody := map[string]string{
		"username":     "adminuser",
		"display_name": "Admin User",
		"password":     "adminpass123",
		"role":         "admin",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	// Make request
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/users", serverBaseURL),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusCreated, resp.StatusCode, string(respBody))
	}
}
