package sqlxdatabase

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/InWamos/trinity-proto/config"
	"github.com/jmoiron/sqlx"
)

func buildPostgresDSN(host, user, password, dbname, port, sslmode string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
}

type SQLXDatabase struct {
	engine *sqlx.DB
	logger *slog.Logger
}

func NewSQLXDatabase(config *config.DatabaseConfig, logger *slog.Logger) (*SQLXDatabase, error) {
	dbLogger := logger.With(
		slog.String("component", "database_engine"),
	)
	dbLogger.Debug("The sqlx database engine has been invoked")

	dsn := buildPostgresDSN(
		config.Address,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName,
		strconv.Itoa(config.Port),
		config.DatabaseSslMode,
	)

	engine, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		dbLogger.Error("failed to connect to database", slog.Any("error", err))
		return nil, err
	}

	// Configure connection pool
	engine.SetMaxOpenConns(25)
	engine.SetMaxIdleConns(5)

	dbLogger.Info("database connection established")

	return &SQLXDatabase{
		engine: engine,
		logger: dbLogger,
	}, nil
}

func (db *SQLXDatabase) GetEngine() *sqlx.DB {
	return db.engine
}

func (db *SQLXDatabase) Dispose() error {
	db.logger.Debug("closing database connection")
	if err := db.engine.Close(); err != nil {
		db.logger.Error("failed to close database connection", slog.Any("error", err))
		return err
	}
	db.logger.Info("database connection closed")
	return nil
}
