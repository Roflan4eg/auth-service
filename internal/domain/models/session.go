package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Session struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	CreatedAt        time.Time `json:"created_at"`
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	UserAgent        string    `json:"user_agent"`
	IpAddress        string    `json:"ip_address"`
	//IsRevoked        bool      `json:"is_revoked"`
}

type JWTClaims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type ValidateTokenResponse struct {
	Valid     bool   `json:"valid"`
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
	Error     string `json:"error"`
}
