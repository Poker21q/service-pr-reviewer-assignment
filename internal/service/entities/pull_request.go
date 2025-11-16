package entities

import (
	"time"

	"github.com/google/uuid"
)

// PullRequest представляет сущность Pull Request (PR).
type PullRequest struct {
	// ID уникальный идентификатор PR.
	ID uuid.UUID

	// Name название PR.
	Name string

	// AuthorID идентификатор автора PR.
	AuthorID uuid.UUID

	// Status статус PR: open или merged.
	Status PullRequestStatus

	// CreatedAt время создания
	CreatedAt time.Time

	// MergedAt время мерджа
	MergedAt *time.Time
}

// PullRequestStatus задает статус PR.
type PullRequestStatus string

const (
	// PullRequestStatusOpen PR открыт.
	PullRequestStatusOpen PullRequestStatus = "open"

	// PullRequestStatusMerged PR смержен.
	PullRequestStatusMerged PullRequestStatus = "merged"
)

type PullRequestReviewer struct {
	ReviewerID    uuid.UUID
	PullRequestID uuid.UUID
}
