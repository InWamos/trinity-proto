package user

import "go.uber.org/fx"

func NewUserModuleContainer() fx.Option {
	// Module with all module's dependencies
	return fx.Module(
		"user_module",
		NewUserApplicationContainer(),
		NewUserInfrastructureContainer(),
		NewUserPresentationContainer(),
	)
}
