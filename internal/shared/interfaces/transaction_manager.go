package interfaces

type TransactionManager interface {
	Commit() error
	Rollback() error
}
