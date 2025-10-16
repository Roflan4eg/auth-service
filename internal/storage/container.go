package storage

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
)

type SQLStorage interface {
	DB() any
	Close() error
}
type CacheStorage interface {
	Client() any
	Close() error
}

type Container struct {
	sqlStorage   SQLStorage
	cacheStorage CacheStorage
	logger       *logger.Logger
}

func NewContainer(cfg *config.Config, logger *logger.Logger) (*Container, error) {
	var (
		err   error
		db    SQLStorage
		cache CacheStorage
	)
	switch cfg.App.DBType {
	case "postgres":
		logger.Debug("Initializing postgres storage")
		db, err = NewPostgresStorage(cfg.Postgres)
	case "mysql":
		//db, err = NewMySQLStorage(cfg.Database)
	}
	if err != nil {
		return nil, err
	}
	switch cfg.App.CacheType {
	case "redis":
		logger.Debug("Initializing redis storage")
		cache, err = NewRedisClient(cfg.Redis)
	case "memcached":
		//db, err = NewMemcachedStorage(cfg.Cache)
	}
	if err != nil {
		return nil, err
	}

	storage := &Container{
		sqlStorage:   db,
		cacheStorage: cache,
		logger:       logger,
	}
	return storage, nil
}

func (c *Container) Cache() any {
	return c.cacheStorage.Client()
}

func (c *Container) SQL() any {
	return c.sqlStorage.DB()
}

func (c *Container) Close(ctx context.Context) error {
	c.logger.Debug("Shutting down storage")
	go c.sqlStorage.Close()
	err := c.cacheStorage.Close()
	if err != nil {
		return err
	}
	return nil
}
