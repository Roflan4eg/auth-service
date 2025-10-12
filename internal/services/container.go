package services

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/repository"
)

type Container struct {
	UserService *UserService
	AuthService *AuthService
}

func NewContainer(
	repository *repository.Container,
	cfg *config.Config,
	logger *logger.Logger,
) *Container {
	userService := NewUserService(repository.UserRepo, logger)
	authService := NewAuthService(repository.SessionRepo, userService, cfg.JWTConfig)

	return &Container{UserService: userService, AuthService: authService}
}
