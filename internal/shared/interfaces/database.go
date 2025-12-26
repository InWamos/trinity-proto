package interfaces //nolint:var-naming // meaningful name

type Database interface {
	GetSession() (any, error)
	Dispose() error
}
