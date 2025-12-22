package auth

import (
	"log/slog"

	redisinfrast "github.com/InWamos/trinity-proto/internal/auth/infrastructure"
	redisdatabase "github.com/InWamos/trinity-proto/internal/auth/infrastructure/database"
	"go.uber.org/fx"
)

func NewAuthInfrastructureContainer() fx.Option {
	return fx.Module(
		"auth_infrastructure",
		fx.Provide(
			// Provides Redis database engine
			redisdatabase.NewRedisDatabase,
			// Provides Redis mapper
			func() *redisinfrast.RedisMapper {
				return &redisinfrast.RedisMapper{}
			},
			// Provides session repository with redis backend
			func(redisDb *redisdatabase.RedisDatabase, mapper *redisinfrast.RedisMapper, logger *slog.Logger) redisinfrast.SessionRepository {
				return redisinfrast.NewRedisSessionRepository(redisDb.GetClient(), mapper, logger)
			},
		),
	)
}
