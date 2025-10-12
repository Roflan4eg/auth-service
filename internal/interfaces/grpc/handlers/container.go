package handlers

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/services"
)

type Container struct {
	UserService *UserGRPCHandler
	AuthService *AuthGRPCHandler
}

func NewContainer(
	services *services.Container,
	cfg *config.Config,
	logger *logger.Logger,
) *Container {

	userHandler := NewUserGRPCHandler(services.UserService)
	authHandler := NewAuthGRPCHandler(services.AuthService)

	return &Container{
		UserService: userHandler,
		AuthService: authHandler,
	}
}
