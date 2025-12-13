package database

import (
	"context"
	"log/slog"
)

// GormTransactionFactory creates request-scoped transactions
// We need this as we have no scopes in uber-fx apart from the application scope
type GormTransactionFactory struct {
	db     *GormDatabase
	logger *slog.Logger
}

func NewGormTransactionFactory(db *GormDatabase, logger *slog.Logger) *GormTransactionFactory {
	return &GormTransactionFactory{
		db:     db,
		logger: logger.With(slog.String("component", "transaction_factory")),
	}
}

// NewTransaction creates a NEW transaction for each request
func (f *GormTransactionFactory) NewTransaction(ctx context.Context) *GormTransactionManager {
	tx := f.db.engine.WithContext(ctx).Begin()
	return &GormTransactionManager{
		transaction: tx,
		logger:      f.logger,
	}
}
