// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrOpenDB          = errors.New("failed to open database")
	ErrMigration       = errors.New("failed to run migrations")
	ErrGormOpen        = errors.New("failed to gorm open")
	ErrGetAllQuestions = errors.New("failed to get all questions")
	ErrCreateQuestion  = errors.New("failed to create question")
	ErrGetQuestion     = errors.New("failed to get question")
	ErrDeleteQuestion  = errors.New("failed to delete question")
	ErrCreateAnswer    = errors.New("failed to create answer")
	ErrGetAnswer       = errors.New("failed to get answer")
	ErrDeleteAnswer    = errors.New("failed to delete answer")
)

type PostgresStorage struct {
	db *gorm.DB
}

func New(cfg Config, log *slog.Logger) (*PostgresStorage, error) {
	const op = "storage.postgres.NewStrorage"
	log = log.With(
		slog.String("op", op),
	)

	sqlDB, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrOpenDB, err)
	}
	log.Info("start migrate...", slog.String("path", cfg.MigrationsPath))
	if err := goose.Up(sqlDB, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrMigration, err)
	}
	//DB seed
	log.Info("start seeding...")

	if cfg.Seed {
		if err := SeedTeams(sqlDB, 20); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if err := SeedUsers(sqlDB, 200); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if err := AssignUsersToTeams(sqlDB); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	gormDB, err := gorm.Open(gormpg.New(gormpg.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrGormOpen, err)
	}

	return &PostgresStorage{db: gormDB}, nil
}
