package interfaces

import "context"

type Broker interface {
	GetProducer() (any, error)
	Dispose(ctx context.Context) error
}
