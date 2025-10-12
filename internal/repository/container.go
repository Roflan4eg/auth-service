package repository

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/storage"
)

type Container struct {
	UserRepo    *UserPgeRepo
	SessionRepo *SessionRedisRepo
}

func NewContainer(
	storage *storage.Container,
	cfg *config.Config,
	logger *logger.Logger,
) *Container {
	userRepo := NewUserRepository(storage.Postgres())
	sessionRepo := NewSessionRedisRepo(storage.Redis(), cfg.Redis.TTL)

	return &Container{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}
