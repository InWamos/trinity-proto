package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestDatabaseSchemaExists(t *testing.T) {
	// This should run using TestMain setup
	config := testContainers.GetDatabaseConfig()
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
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if records schema exists
	var schemaExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = 'records')").
		Scan(&schemaExists)
	if err != nil {
		t.Fatalf("failed to check schema: %v", err)
	}

	if !schemaExists {
		t.Fatal("records schema does not exist")
	}
	t.Log("✓ records schema exists")

	// Check if telegram_users table exists
	var tableExists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'records' AND table_name = 'telegram_users')").
		Scan(&tableExists)
	if err != nil {
		t.Fatalf("failed to check table: %v", err)
	}

	if !tableExists {
		t.Fatal("records.telegram_users table does not exist")
	}
	t.Log("✓ records.telegram_users table exists")

	// List all tables in records schema
	rows, err := db.QueryContext(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'records'")
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}
	defer rows.Close()

	t.Log("Tables in records schema:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Fatalf("failed to scan table name: %v", err)
		}
		t.Logf("  - %s", tableName)
	}
}
