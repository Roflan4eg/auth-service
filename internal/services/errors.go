package services

import "errors"

var (
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrTokenMalformed      = errors.New("token malformed")
	ErrTokenExpired        = errors.New("token expired")
)
