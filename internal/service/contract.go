package service

import (
	"context"
	"time"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/google/uuid"
)

type storage interface {
	// Teams
	CreateTeam(ctx context.Context, teamName string) (*entities.Team, error)
	IsTeamExists(ctx context.Context, teamName string) (bool, error)

	// Users
	CreateOrUpdateUsers(ctx context.Context, users []entities.User) ([]entities.User, error)
	GetUsersByTeamName(ctx context.Context, teamName string) ([]entities.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	IsUserExists(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)

	// PullRequests
	CreatePullRequest(ctx context.Context, pullRequest *entities.PullRequest) (*entities.PullRequest, error)
	CreatePullRequestReviewers(ctx context.Context, pullRequestID uuid.UUID, reviewerIDs []uuid.UUID) error
	GetPullRequestsByReviewerID(ctx context.Context, reviewerID uuid.UUID) ([]entities.PullRequest, error)
	GetPullRequestByID(ctx context.Context, id uuid.UUID) (*entities.PullRequest, error)
	GetPullRequestReviewerIDs(ctx context.Context, pullRequestID uuid.UUID) ([]uuid.UUID, error)
	UpdatePullRequest(ctx context.Context, pullRequest *entities.PullRequest) (*entities.PullRequest, error)
	DeletePullRequestReviewersByReviewerID(ctx context.Context, reviewerID uuid.UUID) error
	DeletePullRequestReviewerByPullRequestIDAndReviewerID(ctx context.Context, pullRequestID uuid.UUID, reviewerID uuid.UUID) error
}

type txManager interface {
	Read(ctx context.Context, fn func(ctx context.Context) error) error
	Write(ctx context.Context, fn func(ctx context.Context) error) error
}

// timeNowFunc для гибкой подмены в тестах
var timeNowFunc = time.Now
