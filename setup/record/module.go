package record

import "go.uber.org/fx"

func NewRecordModuleContainer() fx.Option {
	// Module with all auth module's dependencies
	return fx.Module(
		"record_module",
	)
}
