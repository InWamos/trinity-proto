package application

type Interactor[InputDTO any, OutputDTO any] interface {
	Execute(input InputDTO) (OutputDTO, error)
}
