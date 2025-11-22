package pr

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
)

type Service interface {
	PullRequestCreate(ctx context.Context, pr PullRequest) (PullRequest, error)
	PullRequestMerge(ctx context.Context, id string) (PullRequest, error)
	PullRequestReassign(ctx context.Context, r PostPullRequestReassign) (PullRequest, error)
}

type service struct {
	storage Storage
	log     *slog.Logger
}

func NewService(storage Storage, log *slog.Logger) Service {
	return &service{storage: storage, log: log}
}

func (s *service) PullRequestCreate(ctx context.Context, pr PullRequest) (PullRequest, error) {
	const op = "service.pullRequest.Create"

	teamName, err := s.storage.GetAuthorTeam(pr.AuthorId)
	if err != nil {
		return PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Info("author team", slog.String("TEAM NAME", teamName))

	freeUsers, err := s.storage.GetFreeReviewers(teamName, pr.AuthorId)
	if err != nil {
		return PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	const maxReviewers = 2
	reviewers := make([]string, 0, maxReviewers)
	if len(freeUsers) > maxReviewers {
		rand.Shuffle(len(freeUsers), func(i, j int) {
			freeUsers[i], freeUsers[j] = freeUsers[j], freeUsers[i]
		})
		freeUsers = freeUsers[:maxReviewers]
	}
	for i, user := range freeUsers {
		if i >= maxReviewers {
			break
		}
		reviewers = append(reviewers, user.UserId)

		s.log.Info("reviewer selected", slog.String("USER ID", user.UserId), slog.String("USERNAME", user.Username))
	}

	if len(reviewers) == 0 {
		return PullRequest{}, fmt.Errorf("%s: no available reviewers in team", op)
	}

	newPullRequest := PullRequest{
		AssignedReviewers: reviewers,
		AuthorId:          pr.AuthorId,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            pr.Status,
	}

	if err := s.storage.PullRequestCreate(newPullRequest); err != nil {
		return PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("pull request created", slog.String("pr_id", pr.PullRequestId), slog.Int("reviewers_count", len(reviewers)))

	return newPullRequest, nil
}
func (s *service) PullRequestMerge(ctx context.Context, id string) (PullRequest, error) {
	const op = "service.pullRrquest.Merge"

	pr, err := s.storage.PullRequestMerge(id)
	if err != nil {
		return PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}
	return pr, nil
}

func (s *service) PullRequestReassign(ctx context.Context, r PostPullRequestReassign) (PullRequest, error) {
const op = "service.pull_request.Reassign"

	prResp, err := s.storage.PullRequestReassign(r)
	if err != nil {
		return PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	return prResp, nil
}
