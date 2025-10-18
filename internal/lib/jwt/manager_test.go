package jwt

import (
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.JWTConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid config",
			config: &config.JWTConfig{
				Secret:          "secret",
				AccessTokenTTL:  15 * time.Minute,
				RefreshTokenTTL: 24 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "Zero access token ttl",
			config: &config.JWTConfig{
				Secret:          "secret",
				AccessTokenTTL:  0,
				RefreshTokenTTL: 24 * time.Hour,
			},
			wantErr:     true,
			errContains: "expiration must be positive",
		},
		{
			name: "Zero refresh token ttl",
			config: &config.JWTConfig{
				Secret:          "secret",
				AccessTokenTTL:  15 * time.Minute,
				RefreshTokenTTL: 0,
			},
			wantErr:     true,
			errContains: "expiration must be positive",
		},
		{
			name: "Access ttl greater then refresh ttl",
			config: &config.JWTConfig{
				Secret:          "secret",
				AccessTokenTTL:  24 * time.Hour,
				RefreshTokenTTL: 15 * time.Minute,
			},
			wantErr:     true,
			errContains: "access expiration must be less than refresh",
		},
		{
			name: "Negative ttl",
			config: &config.JWTConfig{
				Secret:          "secret",
				AccessTokenTTL:  -15 * time.Minute,
				RefreshTokenTTL: 24 * time.Hour,
			},
			wantErr:     true,
			errContains: "expiration must be positive",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			manager, err := NewManager(tc.config)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errContains)
				assert.Nil(t, manager)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, manager)
				assert.Equal(t, tc.config, manager.conf)
			}
		})
	}
}

func TestManager_GenerateAccessToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-secret-key-123",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	manager, err := NewManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)

	uid := "user123"
	sesId := "ses123"

	token, err := manager.GenerateAccessToken(uid, sesId)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uid, claims.UserID)
	assert.Equal(t, sesId, claims.SessionID)
	assert.WithinDuration(t, time.Now().Add(cfg.AccessTokenTTL), claims.ExpiresAt.Time, time.Second)
}

func TestManager_GenerateRefreshToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-secret-key-123",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	manager, err := NewManager(cfg)
	require.NoError(t, err)
	require.NotNil(t, manager)

	uid := "user123"
	sesId := "ses123"

	token, err := manager.GenerateRefreshToken(uid, sesId)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uid, claims.UserID)
	assert.Equal(t, sesId, claims.SessionID)
	assert.WithinDuration(t, time.Now().Add(cfg.RefreshTokenTTL), claims.ExpiresAt.Time, time.Second)
}

func TestManager_ValidateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-secret-key-123",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}

	manager, err := NewManager(cfg)
	require.NoError(t, err)

	userID := "user-123"
	sessionID := "session-456"

	token, err := manager.GenerateAccessToken(userID, sessionID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:        "invalid token format",
			token:       "invalid-token-format",
			wantErr:     true,
			errContains: "token malformed",
		},
		{
			name:        "empty token",
			token:       "",
			wantErr:     true,
			errContains: "token malformed",
		},
		{
			name:        "tampered token",
			token:       token[:len(token)-10] + "tampered",
			wantErr:     true,
			errContains: "token malformed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := manager.ValidateToken(tc.token)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errContains)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.UserID)
				assert.Equal(t, sessionID, claims.SessionID)
			}
		})
	}

	t.Run("different secret", func(t *testing.T) {
		otherManager, err := NewManager(&config.JWTConfig{
			Secret:          "different-secret",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		})
		require.NoError(t, err)

		claims, err := otherManager.ValidateToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestManager_TokenExpiration(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:          "test-secret",
		AccessTokenTTL:  500 * time.Millisecond,
		RefreshTokenTTL: time.Second,
	}

	manager, err := NewManager(cfg)
	require.NoError(t, err)

	token, err := manager.GenerateAccessToken("user-123", "session-456")
	require.NoError(t, err)

	t.Run("token valid initially", func(t *testing.T) {
		claims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
	})

	t.Run("token expired", func(t *testing.T) {
		time.Sleep(550 * time.Millisecond)

		claims, err := manager.ValidateToken(token)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "token expired")
		assert.Nil(t, claims)
	})
}
