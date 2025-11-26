package infrastructure

import (
	"context"
	"fmt"
	"strconv"

	"github.com/InWamos/trinity-proto/config"
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

func NewGormDatabase(config config.DatabaseConfig) (*GormDatabase, error) {
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
	})
	if err != nil {
		return nil, err
	}

	return &GormDatabase{engine: engine}, nil
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
