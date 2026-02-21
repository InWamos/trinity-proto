package interfaces

import "context"

type Broker interface {
	GetSyncUserProducer() (any, error)
	Dispose(ctx context.Context) error
}
