package database

import "github.com/InWamos/trinity-proto/internal/shared/interfaces"

type GormTransactionManager struct {
	session *GormSession
}

func NewTransactionManager(session *GormSession) interfaces.TransactionManager {
	return &GormTransactionManager{session: session}
}

func (tm *GormTransactionManager) Commit() error {
	return tm.session.tx.Commit().Error
}

func (tm *GormTransactionManager) Rollback() error {
	return tm.session.tx.Rollback().Error
}
