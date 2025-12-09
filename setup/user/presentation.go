package user

import (
	"github.com/InWamos/trinity-proto/internal/user/presentation/service"
	v1 "github.com/InWamos/trinity-proto/internal/user/presentation/v1"
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
			// Provides User v1 api mux
			v1.NewUserMuxV1,
		),
	)
}

func NewValidator() *validator.Validate {
	return validator.New()
}
