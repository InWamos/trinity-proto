package database

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"gorm.io/gorm"
)

type GormTransactionManager struct {
	transaction *gorm.DB
	logger      *slog.Logger
}

func NewGormTransactionManager(transaction *gorm.DB, logger *slog.Logger) interfaces.TransactionManager {
	gtmLogger := logger.With("component", "gorm_transaction_manager")
	return &GormTransactionManager{transaction: transaction, logger: gtmLogger}
}

func (tm *GormTransactionManager) Commit(ctx context.Context) error {
	tm.logger.DebugContext(ctx, "Commiting transaction")
	if err := tm.transaction.Commit().Error; err != nil {
		tm.logger.Error("Failed to commit transaction", "err", err)
		return err
	}
	tm.logger.DebugContext(ctx, "Transaction commited")
	return nil
}

func (tm *GormTransactionManager) Rollback(ctx context.Context) error {
	tm.logger.DebugContext(ctx, "rolling back transaction")

	if err := tm.transaction.Rollback().Error; err != nil {
		tm.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", err))
		return err
	}

	tm.logger.DebugContext(ctx, "transaction rolled back")
	return nil
}
