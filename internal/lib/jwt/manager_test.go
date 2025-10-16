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
