package storage

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	postgres *PostgresStorage
	redis    *RedisClient
	logger   *logger.Logger
}

func NewStorage(cfg *config.Config, logger *logger.Logger) (*Container, error) {
	logger.Debug("Initializing postgres storage")
	pgStorage, err := NewPostgresStorage(cfg.Database)
	if err != nil {
		return nil, err
	}

	logger.Debug("Initializing redis storage")
	redisClient, err := NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, err
	}

	storage := &Container{
		postgres: pgStorage,
		redis:    redisClient,
		logger:   logger,
	}
	return storage, nil
}

func (c *Container) Redis() *redis.Client {
	return c.redis.Client
}

func (c *Container) Postgres() *pgxpool.Pool {
	return c.postgres.DB()
}

func (c *Container) Close(ctx context.Context) error {
	c.logger.Debug("Shutting down storage")
	go c.postgres.Close()
	err := c.redis.Close()
	if err != nil {
		return err
	}
	return nil
}
