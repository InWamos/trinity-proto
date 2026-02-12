package sqlxdatabase

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
	"github.com/jmoiron/sqlx"
)

type SQLXTransactionManager struct {
	transaction *sqlx.Tx
	logger      *slog.Logger
}

func NewSQLXTransactionManager(transaction *sqlx.Tx, logger *slog.Logger) interfaces.TransactionManager {
	stmLogger := logger.With("component", "sqlx_transaction_manager")
	return &SQLXTransactionManager{transaction: transaction, logger: stmLogger}
}

func (tm *SQLXTransactionManager) Commit(ctx context.Context) error {
	tm.logger.DebugContext(ctx, "committing transaction")
	if err := tm.transaction.Commit(); err != nil {
		tm.logger.ErrorContext(ctx, "failed to commit transaction", slog.Any("error", err))
		return err
	}
	tm.logger.DebugContext(ctx, "transaction committed")
	return nil
}

func (tm *SQLXTransactionManager) Rollback(ctx context.Context) error {
	tm.logger.DebugContext(ctx, "rolling back transaction")

	if err := tm.transaction.Rollback(); err != nil {
		tm.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", err))
		return err
	}

	tm.logger.DebugContext(ctx, "transaction rolled back")
	return nil
}

func (tm *SQLXTransactionManager) GetTransaction() any {
	return tm.transaction
}
