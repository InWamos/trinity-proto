package interfaces

type Database interface {
	GetSession() (any, error)
	Dispose() error
}
