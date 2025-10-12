package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/internal/domain"
	"github.com/Roflan4eg/auth-serivce/internal/domain/models"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type SessionRedisRepo struct {
	client     *redis.Client
	expiration time.Duration
}

func NewSessionRedisRepo(client *redis.Client, expiration time.Duration) *SessionRedisRepo {
	return &SessionRedisRepo{client: client, expiration: expiration}
}

func (r *SessionRedisRepo) Create(ctx context.Context, session *models.Session) error {
	const op = "repository.SessionRedisRepo.Create"

	ok, err := r.Exists(ctx, session.ID)
	if err != nil {
		return err
	}
	if ok {
		return domain.ErrSessionAlreadyExists
	}

	sessionData := map[string]interface{}{
		"user_id":            session.UserID,
		"access_token":       session.AccessToken,
		"refresh_token":      session.RefreshToken,
		"user_agent":         session.UserAgent,
		"ip_address":         session.IpAddress,
		"created_at":         session.CreatedAt.Unix(),
		"expires_at":         session.ExpiresAt.Unix(),
		"refresh_expires_at": session.RefreshExpiresAt.Unix(),
		//"is_revoked":         session.IsRevoked,
	}
	pipe := r.client.Pipeline()
	pipe.HSet(ctx, "session:"+session.ID, sessionData)
	//pipe.HSet(ctx, "session:"+session.ID, "last_activity", time.Now().Unix())
	pipe.Expire(ctx, "session:"+session.ID, r.expiration)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *SessionRedisRepo) GetById(ctx context.Context, sessionID string) (*models.Session, error) {
	const op = "repository.SessionRedisRepo.GetById"
	data, err := r.client.HGetAll(ctx, "session:"+sessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	if len(data) == 0 {
		return nil, domain.ErrSessionNotFound
	}

	return r.unmarshalSession(sessionID, data)
}

func (r *SessionRedisRepo) Update(ctx context.Context, session *models.Session) error {
	const op = "repository.SessionRedisRepo.Update"

	ok, err := r.Exists(ctx, session.ID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrSessionExpired
	}

	sessionData := map[string]interface{}{
		"user_id":            session.UserID,
		"access_token":       session.AccessToken,
		"refresh_token":      session.RefreshToken,
		"user_agent":         session.UserAgent,
		"ip_address":         session.IpAddress,
		"created_at":         session.CreatedAt.Unix(),
		"expires_at":         session.ExpiresAt.Unix(),
		"refresh_expires_at": session.RefreshExpiresAt.Unix(),
		//"is_revoked":         session.IsRevoked,
	}
	_, err = r.client.HSet(ctx, "session:"+session.ID, sessionData).Result()
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

//func (r *SessionRedisRepo) Delete(ctx context.Context, sessionID string) error {
//	data, res := r.client.Del(ctx, "session:"+sessionID).Result()
//	if res != nil {
//		return fmt.Errorf("failed to delete session: %w", res)
//	}
//	if data == 0 {
//		return ErrSessionNotFound
//	}
//	return nil
//}

func (r *SessionRedisRepo) Revoke(ctx context.Context, sessionID string) error {
	const op = "repository.SessionRedisRepo.Revoke"
	ok, err := r.Exists(ctx, sessionID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrSessionExpired
	}
	err = r.client.Del(ctx, "session:"+sessionID).Err()
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *SessionRedisRepo) Exists(ctx context.Context, sessionID string) (bool, error) {
	const op = "repository.SessionRedisRepo.Exists"
	exists, err := r.client.Exists(ctx, "session:"+sessionID).Result()
	if err != nil {
		return false, fmt.Errorf("%s, %w", op, err)
	}
	return exists > 0, nil
}

func (r *SessionRedisRepo) unmarshalSession(sessionID string, data map[string]string) (*models.Session, error) {
	const op = "repository.SessionRedisRepo.unmarshalSession"
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	expiresAt, err := strconv.ParseInt(data["expires_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	refreshExpiresAt, err := strconv.ParseInt(data["refresh_expires_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	//isRevoked, err := strconv.ParseBool(data["is_revoked"])
	//if err != nil {
	//	return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	//}
	return &models.Session{
		ID:               sessionID,
		UserID:           data["user_id"],
		AccessToken:      data["access_token"],
		RefreshToken:     data["refresh_token"],
		UserAgent:        data["user_agent"],
		IpAddress:        data["ip_address"],
		CreatedAt:        time.Unix(createdAt, 0),
		ExpiresAt:        time.Unix(expiresAt, 0),
		RefreshExpiresAt: time.Unix(refreshExpiresAt, 0),
		//IsRevoked:        isRevoked,
	}, nil
}

//func (r *SessionRedisRepo) UpdateAccessExpiration(ctx context.Context, session *models.Session) error {
//	pipe := r.client.Pipeline()
//	pipe.HSet(ctx, "session:"+session.ID, "expires_at", session.ExpiresAt.Unix())
//	_, err := pipe.Exec(ctx)
//	return err
//}
