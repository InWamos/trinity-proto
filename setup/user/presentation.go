//nolint:revive //meaningful name //nolint: nolintlint //
package user

import (
	"github.com/InWamos/trinity-proto/internal/user/presentation/service"
	v1 "github.com/InWamos/trinity-proto/internal/user/presentation/v1"
	"github.com/InWamos/trinity-proto/internal/user/presentation/v1/handlers"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func NewUserPresentationContainer() fx.Option {
	return fx.Module(
		"user_presentation",
		fx.Provide(
			// Provides validator
			NewValidator,
			// Provides form validator
			service.NewTagValidatePostFormValidator,
			// Provides create user handler
			handlers.NewCreateUserHandler,
			// Provides get user handler
			handlers.NewGetUserHandler,
			// Provides promote user handler
			handlers.NewPromoteUserHandler,
			// Provides demote user handler
			handlers.NewDemoteUserHandler,
			// Provides remove user handler
			handlers.NewRemoveUserHandler,
			// Provides User v1 api mux
			v1.NewUserMuxV1,
		),
	)
}

func NewValidator() *validator.Validate {
	return validator.New()
}
