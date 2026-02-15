package sqlxdatabase

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
)

// SQLXTransactionFactory creates request-scoped transactions
// We need this as we have no scopes in uber-fx apart from the application scope.
type SQLXTransactionFactory struct {
	db     *SQLXDatabase
	logger *slog.Logger
}

func NewSQLXTransactionFactory(db *SQLXDatabase, logger *slog.Logger) interfaces.TransactionManagerFactory {
	return &SQLXTransactionFactory{
		db:     db,
		logger: logger.With(slog.String("component", "transaction_factory")),
	}
}

// NewTransaction creates a NEW transaction for each request.
func (f *SQLXTransactionFactory) NewTransaction(ctx context.Context) (interfaces.TransactionManager, error) {
	tx, err := f.db.engine.BeginTxx(ctx, nil)
	if err != nil {
		f.logger.ErrorContext(ctx, "failed to begin transaction", slog.Any("error", err))
		return nil, err
	}

	return &SQLXTransactionManager{
		transaction: tx,
		logger:      f.logger,
	}, nil
}
