package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(cfg *config.DatabaseConfig, logger *logger.Logger) error {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return fmt.Errorf("create migration connection: %w", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create database driver: %w", err)
	}
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("run migrations: %w", err)
	}
	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			logger.Info("No migrations applied yet")
			return nil
		}
		return fmt.Errorf("get migration version: %w", err)
	}
	logger.Info("Migrations applied successfully",
		logger.Uint64("version", version),
		logger.Bool("dirty", dirty))
	return nil
}
