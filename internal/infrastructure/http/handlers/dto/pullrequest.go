package api_dto

import (
	"pr-service/internal/domain/pr"
	"pr-service/internal/infrastructure/http/openapi"
	"time"
)

var now = time.Now()

type PostPullRequestCreateJSONBody struct {
	AuthorId        string `json:"author_id" validate:"required"`
	PullRequestId   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required,min=3"`
}

type PostPullRequestMergeJSONBody struct {
	PullRequestId string `json:"pull_request_id" validate:"required"`
}

func UsersSetIsActiveToModel(p openapi.PostUsersSetIsActiveJSONBody) pr.UsersSetIsActive {
	return pr.UsersSetIsActive{
		IsActive: p.IsActive,
		UserId:   p.UserId,
	}
}

func GetUserToModel(p openapi.GetUsersGetReviewParams) pr.GetReviewParams {
	return pr.GetReviewParams{
		UserId: p.UserId,
	}
}

func GetTeamToModel(t openapi.GetTeamGetParams) pr.TeamName {
	return pr.TeamName{
		TeamName: t.TeamName,
	}
}

func PostPullRequestReassignToModel(req openapi.PostPullRequestReassignJSONBody) pr.PostPullRequestReassign {
	return pr.PostPullRequestReassign{
		OldUserId:     req.OldUserId,
		PullRequestId: req.PullRequestId,
	}
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

func TeamMapToModel(req openapi.Team) pr.Team {
	members := make([]pr.TeamMember, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, pr.TeamMember{
			IsActive: m.IsActive,
			UserId:   m.UserId,
			Username: m.Username,
		})
	}
	return pr.Team{
		Members:  members,
		TeamName: req.TeamName,
	}
}
