package service

import (
	"context"
	"fmt"
	"unicode"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
)

func (s *Service) GetUserPullRequestReviewRequests(
	ctx context.Context,
	userID uuid.UUID,
) ([]entities.PullRequest, error) {
	var userPRs []entities.PullRequest

	err := s.txManager.Read(ctx, func(ctx context.Context) error {
		exists, err := s.storage.IsUserExists(ctx, userID)
		if err != nil {
			return fmt.Errorf("check user exists: %w", err)
		}
		if !exists {
			return &entities.ErrUserNotFound{UserID: pointer.To(userID)}
		}

		userPRs, err = s.storage.GetPullRequestsByReviewerID(ctx, userID)
		if err != nil {
			return fmt.Errorf("get pull requests: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get user PR review requests: %w", err)
	}

	return userPRs, nil
}

func (s *Service) SetUserActiveStatus(
	ctx context.Context,
	userID uuid.UUID,
	isActive bool,
) (*entities.User, error) {
	var user *entities.User

	err := s.txManager.Write(ctx, func(ctx context.Context) error {
		var err error
		user, err = s.storage.GetUserByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("get user: %w", err)
		}

		if user.IsActive == isActive {
			return nil
		}

		if user.IsActive && !isActive {
			if err := s.storage.DeletePullRequestReviewersByReviewerID(ctx, userID); err != nil {
				return fmt.Errorf("delete PR reviewers: %w", err)
			}
		}

		user.IsActive = isActive
		user, err = s.storage.UpdateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("set user active status: %w", err)
	}

	return user, nil
}

func validateUsername(username string) error {
	if len(username) < 2 {
		return &entities.ErrUserNameValidation{Reason: "username too short"}
	}
	if len(username) > 50 {
		return &entities.ErrUserNameValidation{Reason: "username too long"}
	}

	for _, r := range username {
		if !unicode.IsLetter(r) || r > unicode.MaxASCII {
			return &entities.ErrUserNameValidation{Reason: "username must contain only english letters"}
		}
	}

	return nil
}

func checkDuplicateUserIDs(users []entities.User) error {
	seen := make(map[uuid.UUID]bool)
	var duplicates []uuid.UUID

	for _, user := range users {
		if seen[user.ID] {
			duplicates = append(duplicates, user.ID)
		}
		seen[user.ID] = true
	}

	if len(duplicates) > 0 {
		return &entities.ErrDuplicateUserIDs{IDs: duplicates}
	}

	return nil
}
