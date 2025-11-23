package shared

type TransactionManager interface {
	Commit() error
	Rollback() error
}
