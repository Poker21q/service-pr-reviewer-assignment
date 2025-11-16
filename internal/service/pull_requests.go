package service

import (
	"context"
	"fmt"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
)

func (s *Service) CreatePullRequestAndAssignReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
	pullRequestName string,
	authorID uuid.UUID,
) (*entities.PullRequest, []uuid.UUID, error) {
	if err := validatePullRequestName(pullRequestName); err != nil {
		return nil, nil, fmt.Errorf("invalid pull request name: %w", err)
	}

	var (
		pr          *entities.PullRequest
		reviewerIDs []uuid.UUID
	)

	err := s.txManager.Write(ctx, func(ctx context.Context) error {
		author, err := s.storage.GetUserByID(ctx, authorID)
		if err != nil {
			return fmt.Errorf("get author: %w", err)
		}

		teamMembers, err := s.storage.GetUsersByTeamName(ctx, author.TeamName)
		if err != nil {
			return fmt.Errorf("get author team: %w", err)
		}

		pr = &entities.PullRequest{
			ID:        pullRequestID,
			Name:      pullRequestName,
			AuthorID:  authorID,
			Status:    entities.PullRequestStatusOpen,
			CreatedAt: timeNowFunc(),
		}

		pr, err = s.storage.CreatePullRequest(ctx, pr)
		if err != nil {
			return fmt.Errorf("create pull request: %w", err)
		}

		reviewerIDs = s.selectReviewers(authorID, teamMembers)
		if len(reviewerIDs) == 0 {
			return nil
		}

		if err := s.storage.CreatePullRequestReviewers(ctx, pullRequestID, reviewerIDs); err != nil {
			return fmt.Errorf("assign reviewers: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create PR and assign reviewers: %w", err)
	}

	return pr, reviewerIDs, nil
}

func (s *Service) selectReviewers(authorID uuid.UUID, team []entities.User) []uuid.UUID {
	const maxReviewers = 2

	var reviewers []uuid.UUID
	for _, user := range team {
		if user.ID == authorID || !user.IsActive {
			continue
		}

		reviewers = append(reviewers, user.ID)
		if len(reviewers) == maxReviewers {
			break
		}
	}

	return reviewers
}

func (s *Service) MergePullRequestAndGetReviewers(
	ctx context.Context,
	pullRequestID uuid.UUID,
) (*entities.PullRequest, []uuid.UUID, error) {
	var (
		pullRequest *entities.PullRequest
		reviewerIDs []uuid.UUID
	)

	err := s.txManager.Write(ctx, func(ctx context.Context) error {
		var err error
		pullRequest, err = s.storage.GetPullRequestByID(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		reviewerIDs, err = s.storage.GetPullRequestReviewerIDs(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("get reviewers: %w", err)
		}

		if pullRequest.Status == entities.PullRequestStatusMerged {
			return nil
		}

		pullRequest.MergedAt = pointer.To(timeNowFunc())
		pullRequest, err = s.storage.UpdatePullRequest(ctx, pullRequest)
		if err != nil {
			return fmt.Errorf("update pull request: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("merge pull request: %w", err)
	}

	return pullRequest, reviewerIDs, nil
}

func (s *Service) ReassignReviewer(
	ctx context.Context,
	pullRequestID uuid.UUID,
	oldReviewerID uuid.UUID,
) (*entities.PullRequest, []uuid.UUID, uuid.UUID, error) {
	var (
		pr            *entities.PullRequest
		reviewerIDs   []uuid.UUID
		newReviewerID uuid.UUID
	)

	err := s.txManager.Write(ctx, func(ctx context.Context) error {
		pullRequest, err := s.storage.GetPullRequestByID(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		if pullRequest.Status == entities.PullRequestStatusMerged {
			return &entities.ErrPullRequestAlreadyMerged{ID: pullRequestID}
		}

		oldReviewer, err := s.storage.GetUserByID(ctx, oldReviewerID)
		if err != nil {
			return fmt.Errorf("get old reviewer: %w", err)
		}

		teamMembers, err := s.storage.GetUsersByTeamName(ctx, oldReviewer.TeamName)
		if err != nil {
			return fmt.Errorf("get team members: %w", err)
		}

		currentReviewerIDs, err := s.storage.GetPullRequestReviewerIDs(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("get current reviewers: %w", err)
		}

		if !contains(currentReviewerIDs, oldReviewerID) {
			return &entities.ErrReviewerNotAssigned{ID: oldReviewerID}
		}

		newReviewerID = s.findReplacementCandidate(oldReviewerID, currentReviewerIDs, teamMembers)
		if newReviewerID == uuid.Nil {
			return entities.ErrNoReplacementCandidate
		}

		if err := s.storage.DeletePullRequestReviewerByPullRequestIDAndReviewerID(
			ctx, pullRequestID, oldReviewerID,
		); err != nil {
			return fmt.Errorf("delete old reviewer: %w", err)
		}

		if err := s.storage.CreatePullRequestReviewers(
			ctx, pullRequestID, []uuid.UUID{newReviewerID},
		); err != nil {
			return fmt.Errorf("assign new reviewer: %w", err)
		}

		reviewerIDs = replaceReviewer(currentReviewerIDs, oldReviewerID, newReviewerID)
		pr = pullRequest
		return nil
	})
	if err != nil {
		return nil, nil, uuid.Nil, fmt.Errorf("reassign reviewer: %w", err)
	}

	return pr, reviewerIDs, newReviewerID, nil
}

func contains(ids []uuid.UUID, target uuid.UUID) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

func (s *Service) findReplacementCandidate(
	oldReviewerID uuid.UUID,
	currentReviewerIDs []uuid.UUID,
	teamMembers []entities.User,
) uuid.UUID {
	currentSet := make(map[uuid.UUID]bool)
	for _, id := range currentReviewerIDs {
		currentSet[id] = true
	}

	for _, user := range teamMembers {
		if user.ID != oldReviewerID && user.IsActive && !currentSet[user.ID] {
			return user.ID
		}
	}

	return uuid.Nil
}

func replaceReviewer(reviewerIDs []uuid.UUID, oldID, newID uuid.UUID) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(reviewerIDs))
	for _, id := range reviewerIDs {
		if id != oldID {
			result = append(result, id)
		}
	}
	return append(result, newID)
}

func validatePullRequestName(pullRequestName string) error {
	if len(pullRequestName) < 2 {
		return &entities.ErrPullRequestNameValidation{Reason: "team name too short"}
	}
	if len(pullRequestName) > 100 {
		return &entities.ErrPullRequestNameValidation{Reason: "team name too long"}
	}

	return nil
}
