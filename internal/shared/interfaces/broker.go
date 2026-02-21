package interfaces

import "context"

type Broker interface {
	GetSyncUserProducer() any
	Dispose(ctx context.Context) error
}
