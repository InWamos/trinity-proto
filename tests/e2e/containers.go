package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// Match docker-compose.yml versions
	PostgresImage = "postgres:18.1-trixie"
	RedisImage    = "redis:8.4-bookworm"
	MigrateImage  = "migrate/migrate:v4.19.1"

	// Test database credentials
	TestDBName     = "trinity_test"
	TestDBUser     = "trinity"
	TestDBPassword = "testpassword"
)

// TestContainers holds all container instances for E2E tests
type TestContainers struct {
	PostgresContainer *postgres.PostgresContainer
	RedisContainer    *redis.RedisContainer
	Network           testcontainers.Network

	PostgresHost string
	PostgresPort string
	RedisHost    string
	RedisPort    string
}

// SetupContainers initializes all required containers for E2E testing
func SetupContainers(ctx context.Context) (*TestContainers, error) {
	// Create a shared network
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:   "trinity-test-net",
			Driver: "bridge",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	// Start Postgres container
	pgContainer, err := postgres.Run(ctx,
		PostgresImage,
		postgres.WithDatabase(TestDBName),
		postgres.WithUsername(TestDBUser),
		postgres.WithPassword(TestDBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name:     "postgres",
				Networks: []string{"trinity-test-net"},
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}

	pgHost, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres host: %w", err)
	}

	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres port: %w", err)
	}

	// Run migrations with network configuration
	if err := runMigrations(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Start Redis container
	redisContainer, err := redis.Run(ctx,
		RedisImage,
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Networks: []string{"trinity-test-net"},
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start redis: %w", err)
	}

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redis host: %w", err)
	}

	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		return nil, fmt.Errorf("failed to get redis port: %w", err)
	}

	return &TestContainers{
		PostgresContainer: pgContainer,
		RedisContainer:    redisContainer,
		Network:           network,
		PostgresHost:      pgHost,
		PostgresPort:      pgPort.Port(),
		RedisHost:         redisHost,
		RedisPort:         redisPort.Port(),
	}, nil
}

// runMigrations runs database migrations using the migrate container
func runMigrations(ctx context.Context) error {
	// Get the project root directory
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")
	migrationsPath := filepath.Join(projectRoot, "internal", "user", "infrastructure", "migrations")

	// Use container name for internal network communication
	// PostgreSQL container is named "postgres" on the network
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@postgres:5432/%s?sslmode=disable",
		TestDBUser, TestDBPassword, TestDBName,
	)

	migrateContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: MigrateImage,
			Cmd: []string{
				"-path=/migrations",
				fmt.Sprintf("-database=%s", dbURL),
				"up",
			},
			Mounts: testcontainers.Mounts(
				testcontainers.BindMount(migrationsPath, "/migrations"),
			),
			Networks:   []string{"trinity-test-net"},
			WaitingFor: wait.ForExit().WithExitTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("failed to run migrate container: %w", err)
	}
	defer migrateContainer.Terminate(ctx)

	// Check exit code
	state, err := migrateContainer.State(ctx)
	if err != nil {
		return fmt.Errorf("failed to get migrate container state: %w", err)
	}

	if state.ExitCode != 0 {
		logs, _ := migrateContainer.Logs(ctx)
		return fmt.Errorf("migrations failed with exit code %d, logs: %v", state.ExitCode, logs)
	}

	return nil
}

// Teardown terminates all containers and cleans up resources
func (tc *TestContainers) Teardown(ctx context.Context) error {
	var errs []error

	if tc.PostgresContainer != nil {
		if err := tc.PostgresContainer.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate postgres: %w", err))
		}
	}

	if tc.RedisContainer != nil {
		if err := tc.RedisContainer.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate redis: %w", err))
		}
	}

	if tc.Network != nil {
		if err := tc.Network.Remove(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove network: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("teardown errors: %v", errs)
	}

	return nil
}

// GetDatabaseConfig returns environment variables for the test database
func (tc *TestContainers) GetDatabaseConfig() map[string]string {
	return map[string]string{
		"DATABASE_ADDRESS":  tc.PostgresHost,
		"DATABASE_PORT":     tc.PostgresPort,
		"DATABASE_NAME":     TestDBName,
		"DATABASE_USER":     TestDBUser,
		"DATABASE_PASSWORD": TestDBPassword,
		"DATABASE_SSL_MODE": "disable",
	}
}

// GetRedisConfig returns environment variables for the test Redis
func (tc *TestContainers) GetRedisConfig() map[string]string {
	return map[string]string{
		"REDIS_ADDRESS":  tc.RedisHost,
		"REDIS_PORT":     tc.RedisPort,
		"REDIS_DB_AUTH":  "0",
		"REDIS_PASSWORD": "",
	}
}
