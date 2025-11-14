package entities

import "github.com/google/uuid"

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

	// ReviewerIDs список назначенных ревьюверов.
	ReviewerIDs []uuid.UUID

	// NeedMoreReviewers флаг того, что нужно больше ревьюверов.
	NeedMoreReviewers bool
}

// PullRequestStatus задает статус PR.
type PullRequestStatus string

const (
	// PullRequestStatusOpen PR открыт.
	PullRequestStatusOpen PullRequestStatus = "open"

	// PullRequestStatusMerged PR смержен.
	PullRequestStatusMerged PullRequestStatus = "merged"
)
