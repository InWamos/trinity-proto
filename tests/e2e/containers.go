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
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// Match docker-compose.yml versions.

	PostgresImage = "postgres:18.1-trixie"
	RedisImage    = "redis:8.4-bookworm"
	MigrateImage  = "migrate/migrate:v4.19.1"

	// Test database credentials.

	TestDBName     = "trinity_test"
	TestDBUser     = "trinity"
	TestDBPassword = "testpassword"
)

// TestContainers holds all container instances for E2E tests.
type TestContainers struct {
	PostgresContainer *postgres.PostgresContainer
	RedisContainer    *redis.RedisContainer
	Network           *testcontainers.DockerNetwork

	PostgresHost string
	PostgresPort string
	RedisHost    string
	RedisPort    string
}

// SetupContainers initializes all required containers for E2E testing.
func SetupContainers(ctx context.Context) (*TestContainers, error) {
	// Create a shared network
	net, err := network.New(ctx, network.WithDriver("bridge"))
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
		network.WithNetwork([]string{"postgres"}, net),
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
	if err := runMigrations(ctx, net); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Start Redis container
	redisContainer, err := redis.Run(ctx,
		RedisImage,
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
		network.WithNetwork([]string{"redis"}, net),
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
		Network:           net,
		PostgresHost:      pgHost,
		PostgresPort:      pgPort.Port(),
		RedisHost:         redisHost,
		RedisPort:         redisPort.Port(),
	}, nil
}

// runMigrations runs database migrations using the migrate container.
func runMigrations(ctx context.Context, net *testcontainers.DockerNetwork) error {
	// Get the project root directory
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")

	// Use container name for internal network communication
	// PostgreSQL container is named "postgres" on the network
	// Each module uses a separate migrations table to track versions independently
	userDBURL := fmt.Sprintf(
		"postgres://%s:%s@postgres:5432/%s?sslmode=disable&x-migrations-table=schema_migrations_user",
		TestDBUser, TestDBPassword, TestDBName,
	)
	recordDBURL := fmt.Sprintf(
		"postgres://%s:%s@postgres:5432/%s?sslmode=disable&x-migrations-table=schema_migrations_record",
		TestDBUser, TestDBPassword, TestDBName,
	)

	// Run user module migrations
	userMigrationsPath := filepath.Join(projectRoot, "internal", "user", "infrastructure", "migrations")
	if err := runModuleMigrations(ctx, net, userMigrationsPath, userDBURL, "user"); err != nil {
		return err
	}

	// Run record module migrations
	recordMigrationsPath := filepath.Join(projectRoot, "internal", "record", "infrastructure", "migrations")
	if err := runModuleMigrations(ctx, net, recordMigrationsPath, recordDBURL, "record"); err != nil {
		return err
	}

	return nil
}

// runModuleMigrations runs migrations for a specific module.
func runModuleMigrations(
	ctx context.Context,
	net *testcontainers.DockerNetwork,
	migrationsPath, dbURL, moduleName string,
) error {
	migrateReq := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: MigrateImage,
			Cmd: []string{
				"-path=/migrations",
				fmt.Sprintf("-database=%s", dbURL),
				"up",
			},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      migrationsPath,
					ContainerFilePath: "/migrations",
				},
			},
			WaitingFor: wait.ForExit().WithExitTimeout(30 * time.Second),
		},
		Started: true,
	}

	// Attach the migrate container to the network
	if err := network.WithNetwork([]string{}, net)(&migrateReq); err != nil {
		return fmt.Errorf("failed to attach %s migrate container to network: %w", moduleName, err)
	}

	migrateContainer, err := testcontainers.GenericContainer(ctx, migrateReq)
	if err != nil {
		return fmt.Errorf("failed to run %s migrate container: %w", moduleName, err)
	}
	defer migrateContainer.Terminate(ctx)

	// Get logs before checking exit code
	logs, logsErr := migrateContainer.Logs(ctx)
	if logsErr != nil {
		fmt.Printf("Warning: failed to get %s migration logs: %v\n", moduleName, logsErr)
	} else {
		fmt.Printf("=== %s migration logs ===\n", moduleName)
		// Read and print logs
		buf := make([]byte, 4096)
		for {
			n, err := logs.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
		fmt.Printf("=== end %s migration logs ===\n", moduleName)
	}

	// Check exit code
	state, err := migrateContainer.State(ctx)
	if err != nil {
		return fmt.Errorf("failed to get %s migrate container state: %w", moduleName, err)
	}

	if state.ExitCode != 0 {
		return fmt.Errorf("%s migrations failed with exit code %d", moduleName, state.ExitCode)
	}

	fmt.Printf("✓ %s module migrations completed successfully\n", moduleName)
	return nil
}

// Teardown terminates all containers and cleans up resources.
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

// GetDatabaseConfig returns environment variables for the test database.
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

// GetRedisConfig returns environment variables for the test Redis.
func (tc *TestContainers) GetRedisConfig() map[string]string {
	return map[string]string{
		"REDIS_ADDRESS":  tc.RedisHost,
		"REDIS_PORT":     tc.RedisPort,
		"REDIS_DB_AUTH":  "0",
		"REDIS_PASSWORD": "",
	}
}
