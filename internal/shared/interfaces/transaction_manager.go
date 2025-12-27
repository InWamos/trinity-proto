package interfaces //nolint:revive //Meaningful package name for interfaces

import "context"

type TransactionManager interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
