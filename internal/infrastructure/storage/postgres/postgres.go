// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"pr-service/internal/domain/pr"
	pgdto "pr-service/internal/infrastructure/storage/postgres/dto"
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
		if err := SeedUsersWithTeams(sqlDB); err != nil {
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

func (p *PostgresStorage) PullRequestMerge(id string) (pr.PullRequest, error) {
	const op = "storage.postgres.PullRequestMerge"

	res := p.db.Model(&pgdto.PullRequest{}).
		Where("pull_request_id = ? AND status = 'OPEN'", id).
		Updates(map[string]any{
			"status":    "MERGED",
			"merged_at": gorm.Expr("NOW()"),
		})

	if res.Error != nil {
		return pr.PullRequest{}, fmt.Errorf("%s: %w", op, res.Error)
	}

	if res.RowsAffected == 0 {
		var prGorm pgdto.PullRequest
		err := p.db.Preload("AssignedReviewers").
			Where("pull_request_id = ?", id).
			First(&prGorm).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pr.PullRequest{}, ErrNotFound
			}
			return pr.PullRequest{}, fmt.Errorf("%s: %w", op, err)
		}

		if prGorm.Status == "MERGED" {
			return prGorm.ToDomain(), nil
		}

		return pr.PullRequest{}, fmt.Errorf("%s: pull request is %s, not OPEN", op, prGorm.Status)
	}

	var prGorm pgdto.PullRequest
	if err := p.db.Preload("AssignedReviewers").
		Where("pull_request_id = ?", id).
		First(&prGorm).Error; err != nil {

		return pr.PullRequest{}, fmt.Errorf("%s: reload after merge: %w", op, err)
	}

	return prGorm.ToDomain(), nil
}

// storage/postgres/pull_request.go

func (p *PostgresStorage) PullRequestReassign(r pr.PostPullRequestReassign) (pr.PullRequest, error) {
	const op = "storage.postgres.PullRequestReassign"

	var prGorm pgdto.PullRequest

	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("AssignedReviewers").
			First(&prGorm, "pull_request_id = ?", r.PullRequestId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if prGorm.Status == "MERGED" {
			return ErrAlreadyMerged
		}

		found := false
		for _, rev := range prGorm.AssignedReviewers {
			if rev.UserID == r.OldUserId {
				found = true
				break
			}
		}
		if !found {
			return ErrReviewerNotInPR
		}

		var author pgdto.User
		if err := tx.Select("team_name").
			First(&author, "user_id = ?", prGorm.AuthorID).Error; err != nil {
			return err
		}

		if author.TeamName == "" {
			return ErrNotAssigned
		}

		var candidate struct{ UserID string }
		err := tx.Raw(`
            SELECT user_id
            FROM users
            WHERE team_name = ?
              AND is_active = true
              AND user_id != ?
              AND user_id != ?
              AND user_id NOT IN (
                SELECT user_id FROM pull_request_reviewers WHERE pull_request_id = ?
              )
            LIMIT 1
        `, author.TeamName, prGorm.AuthorID, r.OldUserId, r.PullRequestId).
			Scan(&candidate).Error

		if err != nil {
			return err
		}
		if candidate.UserID == "" {
			return ErrNoCandidate
		}

		if err := tx.Exec(`
            DELETE FROM pull_request_reviewers
            WHERE pull_request_id = ? AND user_id = ?
        `, r.PullRequestId, r.OldUserId).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
            INSERT INTO pull_request_reviewers (pull_request_id, user_id)
            VALUES (?, ?)
            ON CONFLICT (pull_request_id, user_id) DO NOTHING
        `, r.PullRequestId, candidate.UserID).Error; err != nil {
			return err
		}

		return tx.Preload("AssignedReviewers").
			First(&prGorm, "pull_request_id = ?", r.PullRequestId).Error
	})

	if err != nil {
		return pr.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	return prGorm.ToDomain(), nil
}

func (p *PostgresStorage) TeamAdd(t pr.Team) (pr.Team, error) {
	if t.TeamName == "" {
		return pr.Team{}, fmt.Errorf("team name required")
	}
	if len(t.Members) == 0 {
		return pr.Team{}, ErrNoCandidate
	}

	tx := p.db.Begin()
	if tx.Error != nil {
		return pr.Team{}, tx.Error
	}
	defer tx.Rollback() 

	var exists int64
	if err := tx.Model(&pgdto.TeamModel{}).
		Where("team_name = ?", t.TeamName).
		Count(&exists).Error; err != nil {
		return pr.Team{}, err
	}
	if exists > 0 {
		return pr.Team{}, ErrTeamExists
	}

	if err := tx.Create(&pgdto.TeamModel{TeamName: t.TeamName}).Error; err != nil {
		return pr.Team{}, err
	}

	for _, m := range t.Members {
		if m.UserId == "" {
			return pr.Team{}, ErrNotFound 
		}
		result := tx.Exec(`
			UPDATE users 
			SET team_name = ?, 
			    is_active = ?, 
			    username = ?
			WHERE user_id = ? 
			  AND (team_name IS NULL OR team_name = ?)`, // опционально: можно запретить перепривязку
			t.TeamName, m.IsActive, m.Username, m.UserId, t.TeamName)

		if result.Error != nil {
			return pr.Team{}, result.Error
		}
		if result.RowsAffected == 0 {
			var cnt int64
			if err := tx.Model(&pgdto.UserModel{}).
				Where("user_id = ?", m.UserId).
				Count(&cnt).Error; err != nil {
				return pr.Team{}, err
			}
			if cnt == 0 {
				return pr.Team{}, ErrNotFound
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return pr.Team{}, err
	}

	var members []pr.TeamMember
	if err := p.db.
		Table("users").
		Where("team_name = ?", t.TeamName).
		Select("user_id", "username", "is_active").
		Scan(&members).Error; err != nil {
		return pr.Team{}, err
	}

	return pr.Team{
		TeamName: t.TeamName,
		Members:  members,
	}, nil
}