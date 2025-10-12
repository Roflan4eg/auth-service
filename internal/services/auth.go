package services

import (
	"context"
	"errors"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/domain"
	"github.com/Roflan4eg/auth-serivce/internal/domain/models"
	"github.com/Roflan4eg/auth-serivce/internal/lib/jwt"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/google/uuid"
	"time"
)

type SessionRepo interface {
	Create(ctx context.Context, session *models.Session) error
	GetById(ctx context.Context, sessionID string) (*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	//Delete(ctx context.Context, sessionID string) error
	Revoke(ctx context.Context, sessionID string) error
	Exists(ctx context.Context, sessionID string) (bool, error)
	//UpdateSessionActivity(ctx context.Context, sessionID string) error
}

type UserClient interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	repo       SessionRepo
	userClient UserClient
	jwtManager *jwt.Manager
}

func NewAuthService(repo SessionRepo, userClient UserClient, conf *config.JWTConfig) *AuthService {
	jwtManager, err := jwt.NewManager(conf)
	if err != nil {
		panic(err)
	}
	return &AuthService{repo: repo, userClient: userClient, jwtManager: jwtManager}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.Session, error) {
	ctx = logger.WithData(ctx, map[string]any{"email": email, "password": password})
	newUser, err := s.userClient.CreateUser(ctx, email, password)
	if err != nil {
		return nil, logger.WrapError(ctx, err) //!!!
	}
	ses, err := s.createSession(ctx, newUser)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}

	return ses, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.Session, error) {
	ctx = logger.WithData(ctx, map[string]any{"email": email, "password": password})
	user, err := s.userClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	isValidPass, err := VerifyPassword(password, user.Password)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	if !isValidPass {
		return nil, logger.WrapError(ctx, domain.ErrInvalidPassword)
	}
	ses, err := s.createSession(ctx, user)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	return ses, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	ctx = logger.WithData(ctx, map[string]any{"session_id": sessionID})
	if err := s.repo.Revoke(ctx, sessionID); err != nil {
		return logger.WrapError(ctx, err)
	}
	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	ctx = logger.WithData(ctx, map[string]any{"refresh_token": refreshToken})
	token, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrInvalidTokenFormat):
			err = ErrInvalidRefreshToken
		case errors.Is(err, jwt.ErrTokenMalformed):
			err = ErrTokenMalformed
		}
		return nil, logger.WrapError(ctx, err)
	}
	ses, err := s.repo.GetById(ctx, token.SessionID)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	if ses.RefreshToken != refreshToken {
		return nil, logger.WrapError(ctx, ErrInvalidRefreshToken)
	}
	if ses.RefreshExpiresAt.Before(time.Now()) {
		return nil, logger.WrapError(ctx, ErrTokenExpired)
	}
	newAccessToken, err := s.jwtManager.GenerateAccessToken(ses.UserID, ses.ID)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(ses.UserID, ses.ID)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	ses.AccessToken = newAccessToken
	ses.RefreshToken = newRefreshToken
	ses.ExpiresAt = time.Now().Add(s.jwtManager.GetAccessTokenTTL())

	err = s.repo.Update(ctx, ses)
	if err != nil {
		return nil, logger.WrapError(ctx, err)
	}
	return ses, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, accessToken string) *models.ValidateTokenResponse {

	resp := &models.ValidateTokenResponse{Valid: false, Error: ""}

	token, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	ses, err := s.repo.GetById(ctx, token.SessionID)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	if ses.AccessToken != accessToken {
		resp.Error = ErrInvalidAccessToken.Error()
		return resp
	}
	if ses.ExpiresAt.Before(time.Now()) {
		resp.Error = ErrTokenMalformed.Error()
		return resp
	}
	resp.UserId = ses.UserID
	resp.SessionId = ses.ID
	resp.Valid = true
	return resp
}

func (s *AuthService) createSession(ctx context.Context, user *models.User) (*models.Session, error) {
	sessUuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), sessUuid.String())
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), sessUuid.String())
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:               sessUuid.String(),
		UserID:           user.ID.String(),
		ExpiresAt:        time.Now().Add(s.jwtManager.GetAccessTokenTTL()),
		RefreshExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenTTL()),
		CreatedAt:        time.Now(),
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		UserAgent:        "",
		IpAddress:        "",
		//IsRevoked:    false,
	}

	err = s.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
