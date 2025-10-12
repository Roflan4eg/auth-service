package jwt

import (
	"errors"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidTokenFormat = errors.New("invalid token")
	ErrFailedGen          = errors.New("failed to generate token")
	ErrTokenMalformed     = errors.New("token malformed")
)

type Manager struct {
	conf *config.JWTConfig
}

func NewManager(conf *config.JWTConfig) (*Manager, error) {
	if conf.AccessTokenTTL <= 0 || conf.RefreshTokenTTL <= 0 {
		return nil, errors.New("expiration must be positive")
	}
	if conf.AccessTokenTTL >= conf.RefreshTokenTTL {
		return nil, errors.New("access expiration must be less than refresh")
	}
	return &Manager{conf: conf}, nil
}

func (m *Manager) GenerateAccessToken(userID, sessionID string) (string, error) {

	claims := &Claims{
		SessionID: sessionID,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.conf.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			//NotBefore: jwt.NewNumericDate(time.Now()),
			//Issuer:    "auth-service",
			//Subject:   userID,
			//ID:        sessionID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString([]byte(m.conf.Secret))
	if err != nil {
		return "", ErrFailedGen
	}
	return res, nil
}

func (m *Manager) GenerateRefreshToken(userID, sessionID string) (string, error) {
	claims := &Claims{
		SessionID: sessionID,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.conf.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			//NotBefore: jwt.NewNumericDate(time.Now()),
			//Issuer:    "auth-service",
			//Subject:   userID,
			//ID:        sessionID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.conf.Secret))
}

func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenFormat
		}
		return []byte(m.conf.Secret), nil
	})

	if err != nil {
		return nil, ErrTokenMalformed
	}

	if !token.Valid {
		return nil, ErrInvalidTokenFormat
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidTokenFormat
	}

	return claims, nil
}

func (m *Manager) GetAccessTokenTTL() time.Duration {
	return m.conf.AccessTokenTTL
}

func (m *Manager) GetRefreshTokenTTL() time.Duration {
	return m.conf.RefreshTokenTTL
}
