package repository

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TODO refactor
type Container struct {
	UserRepo    *UserPgeRepo
	SessionRepo *SessionRedisRepo
}

func NewContainer(
	storage *storage.Container,
	cfg *config.Config,
	logger *logger.Logger,
) *Container {
	var (
		userRepo    *UserPgeRepo
		sessionRepo *SessionRedisRepo
	)
	if db, ok := storage.SQL().(*pgxpool.Pool); ok {
		userRepo = NewPgUserRepository(db)
	}

	if cache, ok := storage.Cache().(*redis.Client); ok {
		sessionRepo = NewSessionRedisRepo(cache, cfg.Redis.TTL)
	}

	return &Container{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}
