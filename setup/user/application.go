package user

import (
	"github.com/InWamos/trinity-proto/internal/user/application"
	"github.com/InWamos/trinity-proto/internal/user/application/service"
	"go.uber.org/fx"
)

func NewUserApplicationContainer() fx.Option {
	return fx.Module(
		"user_presentation",
		fx.Provide(
			// Provides password hasher
			service.NewBcryptPasswordHasher,
			// Provides uuid generator
			service.NewUUIDGenerator,
			// Provides CreateUserInteractor
			application.NewCreateUser,
		),
	)
}
