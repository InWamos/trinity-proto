package database

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/InWamos/trinity-proto/config"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func buildGormDSN(host, user, password, dbname, port, sslmode string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
}

type GormDatabase struct {
	engine *gorm.DB
}

func NewGormDatabase(config *config.DatabaseConfig, logger *slog.Logger) (*GormDatabase, error) {
	gdlogger := logger.With(
		slog.String("component", "database_engine"),
	)
	gdlogger.Debug("The Gorm database engine has been invoked")
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(gdlogger.Handler()),
		slogGorm.WithTraceAll(),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelInfo),
	)
	dsn := buildGormDSN(
		config.Address,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName,
		strconv.Itoa(config.Port),
		config.DatabaseSslMode,
	)
	// Here never use an auto-commit.
	// Horrible feature, you will loose a fraction of your brain cells by leaving it enabled
	engine, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
	})
	if err != nil {
		return nil, err
	}

	return &GormDatabase{engine: engine}, nil
}

func (db *GormDatabase) GetEngine() *gorm.DB {
	return db.engine
}

func (db *GormDatabase) GetSession(ctx context.Context) *GormSession {
	tx := db.engine.WithContext(ctx).Begin()
	return &GormSession{tx: tx}
}

func (db *GormDatabase) Dispose() error {
	sqlDB, err := db.engine.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
