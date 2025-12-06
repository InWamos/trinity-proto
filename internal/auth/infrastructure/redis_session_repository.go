package infrastructure

import (
	"context"
	"log/slog"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisSessionRepository struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewRedisSessionRepository(redisClient *redis.Client, logger *slog.Logger) SessionRepository {
	return &RedisSessionRepository{redisClient: redisClient, logger: logger}
}

func (repo *RedisSessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (domain.Session, error)
func (repo *RedisSessionRepository) RevokeSessionByID(ctx context.Context, id uuid.UUID) error

func (repo *RedisSessionRepository) GetAllSessionsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Session, error)
