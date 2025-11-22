// pgdto/pull_request.go
package pgdto

import (
	"pr-service/internal/domain/pr"
	"time"
)

type PullRequest struct {
    PullRequestID   string    `gorm:"column:pull_request_id;primaryKey;type:text"`
    PullRequestName string    `gorm:"column:pull_request_name;type:text;not null"`
    AuthorID        string    `gorm:"column:author_id;type:text;not null"`
    Status          string    `gorm:"column:status;type:text;not null"`
    CreatedAt       time.Time `gorm:"column:created_at"`
    MergedAt        *time.Time `gorm:"column:merged_at"`

    AssignedReviewers []User `gorm:"many2many:pull_request_reviewers;foreignKey:PullRequestID;joinForeignKey:pull_request_id;References:UserID;joinReferences:user_id"`
}

type User struct {
    UserID string `gorm:"column:user_id;primaryKey;type:text"`
}

type PullRequestReviewer struct {
	PullRequestID string `gorm:"primaryKey;column:pull_request_id;type:text"`
	UserID        string `gorm:"primaryKey;column:user_id;type:text"`
}

func (p *PullRequest) ToDomain() pr.PullRequest {
	reviewers := make([]string, 0, len(p.AssignedReviewers))
	for _, u := range p.AssignedReviewers {
		reviewers = append(reviewers, u.UserID)
	}

	return pr.PullRequest{
		PullRequestId:     p.PullRequestID,
		PullRequestName:   p.PullRequestName,
		AuthorId:          p.AuthorID,
		Status:            p.Status,
		CreatedAt:         &p.CreatedAt,
		MergedAt:          p.MergedAt,
		AssignedReviewers: reviewers,
	}
}