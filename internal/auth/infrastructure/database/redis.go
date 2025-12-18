package database

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/InWamos/trinity-proto/config"
	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisDatabase(config *config.RedisConfig, logger *slog.Logger) (*RedisDatabase, error) {
	redisLogger := logger.With(slog.String("component", "redis_engine"))
	redisLogger.Debug("The Redis database engine has been invoked")

	dbNumber, err := strconv.Atoi(config.DbNumberAuth)
	if err != nil {
		redisLogger.Error("invalid Redis DB number", slog.String("db", config.DbNumberAuth))
		return nil, fmt.Errorf("invalid Redis DB number: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           dbNumber,
		Protocol:     3, // Use RESP3 protocol
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		redisLogger.Error("failed to connect to Redis", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisLogger.Info("successfully connected to Redis",
		slog.String("host", config.Host),
		slog.String("port", config.Port),
		slog.Int("db", dbNumber),
		slog.Int("protocol", 3),
	)

	return &RedisDatabase{
		client: client,
		logger: redisLogger,
	}, nil
}

func (rd *RedisDatabase) GetClient() *redis.Client {
	return rd.client
}

func (rd *RedisDatabase) Close() error {
	rd.logger.Info("closing Redis connection")
	return rd.client.Close()
}

func (rd *RedisDatabase) Ping(ctx context.Context) error {
	return rd.client.Ping(ctx).Err()
}

func (rd *RedisDatabase) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := rd.Ping(ctx); err != nil {
		rd.logger.Error("Redis health check failed", slog.String("err", err.Error()))
		return err
	}

	rd.logger.Debug("Redis health check passed")
	return nil
}
