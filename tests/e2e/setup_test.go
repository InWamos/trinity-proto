package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	testContainers *TestContainers
	serverBaseURL  string
)

func TestMain(m *testing.M) {
	fmt.Println("=== Starting E2E Tests ===")
	fmt.Println("Setting up test containers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Setup containers
	var err error
	fmt.Println("Creating PostgreSQL and Redis containers...")
	testContainers, err = SetupContainers(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup containers: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Containers started successfully")

	// Set environment variables for the test
	fmt.Println("Setting environment variables...")
	for k, v := range testContainers.GetDatabaseConfig() {
		os.Setenv(k, v)
		fmt.Printf("  %s=%s\n", k, v)
	}
	for k, v := range testContainers.GetRedisConfig() {
		os.Setenv(k, v)
		fmt.Printf("  %s=%s\n", k, v)
	}

	// Server config
	os.Setenv("SERVER_ADDRESS", "127.0.0.1")
	os.Setenv("SERVER_PORT", "18080") // Use different port for tests
	os.Setenv("SERVER_TRUSTED_PROXY", "127.0.0.1")
	os.Setenv("SERVER_ALLOWED_ORIGIN", "*")
	os.Setenv("LOGGING_LEVEL", "debug")

	serverBaseURL = "http://127.0.0.1:18080"
	fmt.Printf("✓ Server will run on %s\n", serverBaseURL)

	// Run tests
	fmt.Println("=== Running Tests ===")
	code := m.Run()

	// Teardown
	fmt.Println("=== Tearing Down ===")
	if err := testContainers.Teardown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to teardown containers: %v\n", err)
	}
	fmt.Println("✓ Teardown complete")

	os.Exit(code)
}
