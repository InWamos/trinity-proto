package interfaces //nolint:var-naming // meaningful name

import "context"

type Interactor[InputDTO any, OutputDTO any] interface {
	Execute(ctx context.Context, input InputDTO) (OutputDTO, error)
}
