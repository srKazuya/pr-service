package api_dto

import (
	"pr-service/internal/domain/pr"
	"time"
)

var now = time.Now()

type PostPullRequestCreateJSONBody struct {
	AuthorId        string `json:"author_id" validate:"required"`
	PullRequestId   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required,min=3"`
}


func PostPullRequestMapToModel(req PostPullRequestCreateJSONBody) pr.PullRequest {
	return pr.PullRequest{
		AssignedReviewers: nil,
		AuthorId:          req.AuthorId,
		CreatedAt:         &now,
		MergedAt:          nil,
		PullRequestId:     req.PullRequestId,
		PullRequestName:   req.PullRequestName,
		Status:            "OPEN",
	}
}
