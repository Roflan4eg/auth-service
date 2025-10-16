package storage

import (
	"context"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

func (s *RedisClient) Close() error {
	err := s.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}
	return nil
}

func (s *RedisClient) Client() any {
	return s.client
}
