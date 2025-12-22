package auth

import "go.uber.org/fx"

func NewAuthModuleContainer() fx.Option {
	// Module with all auth module's dependencies
	return fx.Module(
		"auth_module",
		NewAuthApplicationContainer(),
		NewAuthInfrastructureContainer(),
		NewAuthPresentationContainer(),
	)
}
