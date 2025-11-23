package pr

import "time"

type PullRequest struct {
	// AssignedReviewers user_id назначенных ревьюверов (0..2)
	AssignedReviewers []string   
	AuthorId          string    
	CreatedAt         *time.Time 
	MergedAt          *time.Time 
	PullRequestId     string     
	PullRequestName   string     
	Status            string     
}

type User struct {
	IsActive bool   
	TeamName string 
	UserId   string 
	Username string 
}

type PostPullRequestReassign struct {
	OldUserId     string 
	PullRequestId string 
}

type Team struct {
	Members  []TeamMember 
	TeamName string       
}

// TeamMember defines model for TeamMember.
type TeamMember struct {
	IsActive bool   
	UserId   string 
	Username string 
}