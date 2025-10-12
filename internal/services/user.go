package services

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/internal/domain"
	"github.com/Roflan4eg/auth-serivce/internal/domain/models"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	storage UserRepo
	logger  *logger.Logger
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUserPassword(ctx context.Context, user *models.User) error
}

func NewUserService(storage UserRepo, logger *logger.Logger) *UserService {
	return &UserService{storage: storage, logger: logger}
}

func (s *UserService) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	ctx = logger.WithData(ctx, map[string]any{
		"email":    email,
		"password": password,
	})
	pass, err := HashPass(password)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	user := &models.User{
		ID:        id,
		Email:     email,
		CreatedAt: time.Now(),
		IsActive:  true,
		Password:  pass,
	}
	err = s.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, uid string) (*models.User, error) {
	ctx = logger.WithData(ctx, map[string]any{"uid": uid})
	user, err := s.storage.GetUserByID(ctx, uid)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx = logger.WithData(ctx, map[string]any{"email": email})
	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	return user, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, uid, oldPassword, newPassword string) error {
	ctx = logger.WithData(ctx, map[string]any{"uid": uid, "oldPassword": oldPassword, "newPassword": newPassword})
	user, err := s.storage.GetUserByID(ctx, uid)
	if err != nil {
		return logger.WrapError(ctx, err)
	}

	isValidPass, err := VerifyPassword(oldPassword, user.Password)
	if err != nil {
		return logger.WrapError(ctx, err)
	}
	if !isValidPass {
		return logger.WrapError(ctx, domain.ErrInvalidPassword)
	}
	newPass, err := HashPass(newPassword)
	if err != nil {
		return err
	}
	user.Password = newPass

	if err = s.storage.UpdateUserPassword(ctx, user); err != nil {
		return logger.WrapError(ctx, err)
	}
	return nil
}
