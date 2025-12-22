package interfaces

import "context"

type TransactionManager interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
