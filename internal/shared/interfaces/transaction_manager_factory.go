package interfaces

import "context"

type TransactionManagerFactory interface {
	NewTransaction(ctx context.Context) (TransactionManager, error)
}
