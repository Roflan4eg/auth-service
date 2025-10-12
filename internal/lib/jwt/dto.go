package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID    string `json:"uid"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}
