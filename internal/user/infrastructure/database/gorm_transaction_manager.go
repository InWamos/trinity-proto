package database

import (
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/shared/interfaces"
)

type GormTransactionManager struct {
	session *GormSession
	logger  *slog.Logger
}

func NewGormTransactionManager(session *GormSession, logger *slog.Logger) interfaces.TransactionManager {
	gtmLogger := logger.With("component", "gorm_transaction_manager")
	return &GormTransactionManager{session: session, logger: gtmLogger}
}

func (tm *GormTransactionManager) Commit() error {
	tm.logger.Debug("Commiting transaction")
	if err := tm.session.tx.Commit().Error; err != nil {
		tm.logger.Error("Failed to commit transaction", "err", err)
		return err
	}
	tm.logger.Debug("Transaction commited")
	return nil
}

func (tm *GormTransactionManager) Rollback() error {
	return tm.session.tx.Rollback().Error
}
