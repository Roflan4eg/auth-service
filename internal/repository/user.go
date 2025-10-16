package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/internal/domain"
	"github.com/Roflan4eg/auth-serivce/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPgeRepo struct {
	db *pgxpool.Pool
}

func NewPgUserRepository(db *pgxpool.Pool) *UserPgeRepo {
	return &UserPgeRepo{db: db}
}

func (r *UserPgeRepo) CreateUser(ctx context.Context, user *models.User) error {
	const op = "repository.UserPgeRepo.CreateUser"
	query := `INSERT INTO users (id, email, password, created_at, is_active) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password, user.CreatedAt, user.IsActive)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *UserPgeRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "repository.UserPgeRepo.GetUserByEmail"
	var user models.User
	query := `SELECT id, email, password, created_at, is_active FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &user, nil
}

func (r *UserPgeRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const op = "repository.UserPgeRepo.GetUserByID"
	var user models.User
	err := r.db.QueryRow(ctx, `SELECT id, email, password, created_at, is_active FROM users WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &user, nil
}

func (r *UserPgeRepo) UpdateUserPassword(ctx context.Context, user *models.User) error {
	const op = "repository.UserPgeRepo.UpdateUserPassword"
	query := `UPDATE users SET password = $1 WHERE id = $2`
	res, err := r.db.Exec(ctx, query, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
