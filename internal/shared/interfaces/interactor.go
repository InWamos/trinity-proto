package interfaces

import "context"

type Interactor[InputDTO any, OutputDTO any] interface {
	Execute(ctx context.Context, input InputDTO) (OutputDTO, error)
}
