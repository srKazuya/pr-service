// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"pr-service/internal/domain/pr"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
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

// PullRequestCreate — создаёт PR + сразу назначает ревьюеров (если переданы)
func (p *PostgresStorage) PullRequestCreate(prEntity pr.PullRequest) error {
	const op = "storage.postgres.PullRequestCreate"

	return p.db.Transaction(func(tx *gorm.DB) error {
		var exists int64
		if err := tx.Table("pull_requests").
			Where("pull_request_id = ?", prEntity.PullRequestId).
			Count(&exists).Error; err != nil {
			return fmt.Errorf("%s: check existence failed: %w", op, err)
		}
		if exists > 0 {
			return ErrPrExists
		}
		now := time.Now()
		if prEntity.CreatedAt == nil {
			prEntity.CreatedAt = &now
		}
		
		prToInsert := map[string]interface{}{
			"pull_request_id":   prEntity.PullRequestId,
			"pull_request_name": prEntity.PullRequestName,
			"author_id":         prEntity.AuthorId,
			"status":            prEntity.Status,
			"created_at":        prEntity.CreatedAt,
			"merged_at":         prEntity.MergedAt,
		}

		if err := tx.Table("pull_requests").Create(prToInsert).Error; err != nil {
			return fmt.Errorf("%s: failed to create pull_request: %w", op, err)
		}

		if len(prEntity.AssignedReviewers) > 0 {
			records := make([]map[string]interface{}, 0, len(prEntity.AssignedReviewers))
			for _, userID := range prEntity.AssignedReviewers {
				records = append(records, map[string]interface{}{
					"pull_request_id": prEntity.PullRequestId,
					"user_id":         userID,
				})
			}

			if err := tx.Table("pull_request_reviewers").
				CreateInBatches(records, 100).Error; err != nil {
				return fmt.Errorf("%s: failed to assign reviewers: %w", op, err)
			}
		}

		return nil 
	})
}

// UsersSetIsActive — активирует пользователя
func (p *PostgresStorage) UsersSetIsActive(userID string) error {
	const op = "storage.postgres.UsersSetIsActive"

	result := p.db.
		Table("users").
		Where("user_id = ?", userID).
		Update("is_active", true)

	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetAuthorTeam — возвращает team_name автора
func (p *PostgresStorage) GetAuthorTeam(userID string) (string, error) {
	const op = "storage.postgres.GetAuthorTeam"
	var teamName string
	err := p.db.
		Table("users").
		Select("team_name").
		Where("user_id = ?", userID).
		Scan(&teamName).Error
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return teamName, nil
}

// GetFreeReviewers — возвращает свободных активных ревьюеров из той же команды
func (p *PostgresStorage) GetFreeReviewers(teamName string, authorUserID string) ([]pr.User, error) {
	const op = "storage.postgres.GetFreeReviewers"

	var users []pr.User

	err := p.db.
		Table("users").
		Select(`
			user_id AS user_id,
			username AS username,
			team_name AS team_name,
			is_active AS is_active
		`).
		Where("team_name = ? AND user_id != ? AND is_active != ?", teamName, authorUserID, false).
		Scan(&users).Error

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(users) == 0 {
		return nil, ErrNoCandidate
	}

	return users, nil
}
