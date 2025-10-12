package storage

import (
	"context"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(cfg *config.DatabaseConfig) (*PostgresStorage, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStorage{pool: pool}, nil
}

func (s *PostgresStorage) Close() {
	s.pool.Close()
}

func (s *PostgresStorage) DB() *pgxpool.Pool {
	return s.pool
}
