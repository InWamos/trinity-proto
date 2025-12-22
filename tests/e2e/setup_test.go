package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	// Initialize default test users in the database
	if err := initializeTestUsers(testContainers); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize test users: %v\n", err)
		os.Exit(1)
	}

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

// initializeTestUsers creates default test users (admin and regular user) in the database
func initializeTestUsers(tc *TestContainers) error {
	config := tc.GetDatabaseConfig()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config["DATABASE_USER"],
		config["DATABASE_PASSWORD"],
		config["DATABASE_ADDRESS"],
		config["DATABASE_PORT"],
		config["DATABASE_NAME"],
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Hash passwords
	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	userHash, err := bcrypt.GenerateFromPassword([]byte("user12345"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %w", err)
	}

	// Insert admin user
	adminID := uuid.New()
	adminQuery := `
		INSERT INTO "user".users (id, username, display_name, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (username) DO NOTHING
	`
	if _, err := db.Exec(adminQuery, adminID, "admin", "Admin User", string(adminHash), "admin"); err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	// Insert regular user
	userID := uuid.New()
	userQuery := `
		INSERT INTO "user".users (id, username, display_name, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (username) DO NOTHING
	`
	if _, err := db.Exec(userQuery, userID, "testuser", "Test User", string(userHash), "user"); err != nil {
		return fmt.Errorf("failed to insert test user: %w", err)
	}

	return nil
}
