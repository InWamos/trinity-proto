package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/InWamos/trinity-proto/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisSessionRepository struct {
	redisClient *redis.Client
	redisMapper *RedisMapper
	logger      *slog.Logger
}

func NewRedisSessionRepository(
	redisClient *redis.Client,
	redisMapper *RedisMapper,
	logger *slog.Logger,
) SessionRepository {
	return &RedisSessionRepository{redisClient: redisClient, redisMapper: redisMapper, logger: logger}
}

func (repo *RedisSessionRepository) GetSessionByToken(ctx context.Context, token string) (domain.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", token)
	result, err := repo.redisClient.HGetAll(ctx, sessionKey).Result()
	if err != nil {
		repo.logger.ErrorContext(ctx, "failed to get session by token", slog.String("err", err.Error()))
		return domain.Session{}, ErrInternal
	}

	if len(result) == 0 {
		repo.logger.DebugContext(ctx, "session not found", slog.String("token", token))
		return domain.Session{}, ErrSessionNotFound
	}

	data := make(map[string]any)
	for k, v := range result {
		data[k] = v
	}

	return repo.redisMapper.MapToSession(data, token)
}

func (repo *RedisSessionRepository) RevokeSessionByToken(ctx context.Context, token string) error {
	keysAffected, err := repo.redisClient.Del(ctx, "session:"+token).Result()
	if err != nil {
		repo.logger.ErrorContext(ctx, "failed to revoke session by token", slog.String("err", err.Error()))
		return ErrInternal
	}
	repo.logger.DebugContext(ctx, "revoked session", slog.Int64("keys affected", keysAffected))
	return nil
}

func (repo *RedisSessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	data := repo.redisMapper.SessionToMap(session)

	err := repo.redisClient.HSet(ctx, "session:"+session.Token, data).Err()
	if err != nil {
		repo.logger.ErrorContext(ctx, "failed to create session", slog.Any("err", err))
		return err
	}

	ttl := time.Until(session.ExpiresAt)
	repo.redisClient.Expire(ctx, "session:"+session.Token, ttl)
	repo.logger.DebugContext(
		ctx,
		"session has been created",
		slog.String("session_id", data["id"].(string)),
		slog.String("user_id", data["user_id"].(string)),
		slog.String("user_role", data["user_role"].(string)),
	)
	return nil
}

func (repo *RedisSessionRepository) GetAllSessionsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Session, error) {
	pattern := "session:*"
	var sessions []domain.Session

	iter := repo.redisClient.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		// Extract token from key (format: "session:TOKEN")
		token := key[8:] // Skip "session:" prefix

		result, err := repo.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}

		data := make(map[string]interface{})
		for k, v := range result {
			data[k] = v
		}

		session, err := repo.redisMapper.MapToSession(data, token)
		if err != nil {
			continue
		}

		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	if err := iter.Err(); err != nil {
		repo.logger.ErrorContext(ctx, "error scanning sessions", slog.Any("err", err))
		return nil, err
	}

	return sessions, nil
}
