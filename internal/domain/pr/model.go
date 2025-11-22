package pr

import "time"

type PullRequest struct {
	// AssignedReviewers user_id назначенных ревьюверов (0..2)
	AssignedReviewers []string   `gorm:"many2many:pull_request_reviewers;joinForeignKey:PullRequestId;joinReferences:UserId"`
	AuthorId          string     `json:"author_id"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
	PullRequestId     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	Status            string     `json:"status"`
}

type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserId   string `gorm:"primaryKey" json:"user_id" `
	Username string `json:"username"`
}
