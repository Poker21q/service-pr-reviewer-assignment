package storage

import (
	"context"
	"fmt"
	"time"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/google/uuid"
)

type pullRequestDB struct {
	ID        uuid.UUID
	Name      string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	MergedAt  *time.Time
}

func (s *Storage) GetPullRequestsByReviewerID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entities.PullRequest, error) {
	const query = `
		SELECT id, name, author_id, created_at, merged_at
		FROM pull_requests
		WHERE id IN (
			SELECT pull_request_id
			FROM pull_request_reviewers
			WHERE reviewer_id = $1
		)
	`

	rows, err := s.querier.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query pull requests by reviewer: %w", err)
	}
	defer rows.Close()

	var prsDB []pullRequestDB
	for rows.Next() {
		var prDB pullRequestDB
		if err := rows.Scan(&prDB.ID, &prDB.Name, &prDB.AuthorID, &prDB.CreatedAt, &prDB.MergedAt); err != nil {
			return nil, fmt.Errorf("scan pull request: %w", err)
		}
		prsDB = append(prsDB, prDB)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return convertPullRequestsDBToEntities(prsDB), nil
}

func (s *Storage) CreatePullRequest(ctx context.Context, pullRequest *entities.PullRequest) (*entities.PullRequest, error) {
	const query = `
		INSERT INTO pull_requests (id, name, author_id, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, author_id, created_at, merged_at
	`

	var prDB pullRequestDB
	err := s.querier.QueryRow(
		ctx,
		query,
		pullRequest.ID,
		pullRequest.Name,
		pullRequest.AuthorID,
		pullRequest.CreatedAt,
		pullRequest.MergedAt,
	).Scan(&prDB.ID, &prDB.Name, &prDB.AuthorID, &prDB.CreatedAt, &prDB.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	result := convertPullRequestDBToEntity(prDB)
	return &result, nil
}

func (s *Storage) CreatePullRequestReviewers(ctx context.Context, pullRequestID uuid.UUID, reviewerIDs []uuid.UUID) error {
	if len(reviewerIDs) == 0 {
		return nil
	}

	builder := s.stmtBuilder.
		Insert("pull_request_reviewers").
		Columns("pull_request_id", "reviewer_id")

	for _, reviewerID := range reviewerIDs {
		builder = builder.Values(pullRequestID, reviewerID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	_, err = s.querier.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert reviewers: %w", err)
	}

	return nil
}

func (s *Storage) GetPullRequestByID(ctx context.Context, id uuid.UUID) (*entities.PullRequest, error) {
	const query = `SELECT id, name, author_id, created_at, merged_at FROM pull_requests WHERE id = $1`

	var prDB pullRequestDB
	err := s.querier.QueryRow(ctx, query, id).Scan(
		&prDB.ID,
		&prDB.Name,
		&prDB.AuthorID,
		&prDB.CreatedAt,
		&prDB.MergedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get pull request by id: %w", err)
	}

	result := convertPullRequestDBToEntity(prDB)
	return &result, nil
}

func (s *Storage) GetPullRequestReviewerIDs(ctx context.Context, pullRequestID uuid.UUID) ([]uuid.UUID, error) {
	const query = `SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1`

	rows, err := s.querier.Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("query reviewers: %w", err)
	}
	defer rows.Close()

	var reviewerIDs []uuid.UUID
	for rows.Next() {
		var reviewerID uuid.UUID
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("scan reviewer id: %w", err)
		}
		reviewerIDs = append(reviewerIDs, reviewerID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return reviewerIDs, nil
}

func (s *Storage) UpdatePullRequest(ctx context.Context, pullRequest *entities.PullRequest) (*entities.PullRequest, error) {
	const query = `
		UPDATE pull_requests 
		SET name = $2, author_id = $3, merged_at = $4
		WHERE id = $1
		RETURNING id, name, author_id, created_at, merged_at
	`

	var prDB pullRequestDB
	err := s.querier.QueryRow(
		ctx,
		query,
		pullRequest.ID,
		pullRequest.Name,
		pullRequest.AuthorID,
		pullRequest.MergedAt,
	).Scan(&prDB.ID, &prDB.Name, &prDB.AuthorID, &prDB.CreatedAt, &prDB.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("update pull request: %w", err)
	}

	result := convertPullRequestDBToEntity(prDB)
	return &result, nil
}

func (s *Storage) DeletePullRequestReviewersByReviewerID(ctx context.Context, reviewerID uuid.UUID) error {
	const query = `DELETE FROM pull_request_reviewers WHERE reviewer_id = $1`

	_, err := s.querier.Exec(ctx, query, reviewerID)
	if err != nil {
		return fmt.Errorf("delete reviewers by reviewer id: %w", err)
	}

	return nil
}

func (s *Storage) DeletePullRequestReviewerByPullRequestIDAndReviewerID(
	ctx context.Context,
	pullRequestID uuid.UUID,
	reviewerID uuid.UUID,
) error {
	const query = `DELETE FROM pull_request_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2`

	_, err := s.querier.Exec(ctx, query, pullRequestID, reviewerID)
	if err != nil {
		return fmt.Errorf("delete reviewer: %w", err)
	}

	return nil
}

func convertPullRequestsDBToEntities(prsDB []pullRequestDB) []entities.PullRequest {
	out := make([]entities.PullRequest, 0, len(prsDB))
	for _, prDB := range prsDB {
		out = append(out, convertPullRequestDBToEntity(prDB))
	}
	return out
}

func convertPullRequestDBToEntity(prDB pullRequestDB) entities.PullRequest {
	return entities.PullRequest{
		ID:        prDB.ID,
		Name:      prDB.Name,
		AuthorID:  prDB.AuthorID,
		Status:    statusByMergedAt(prDB.MergedAt),
		CreatedAt: prDB.CreatedAt,
		MergedAt:  prDB.MergedAt,
	}
}

func statusByMergedAt(mergedAt *time.Time) entities.PullRequestStatus {
	if mergedAt == nil {
		return entities.PullRequestStatusOpen
	}
	return entities.PullRequestStatusMerged
}
