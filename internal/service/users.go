package service

import (
	"context"
	"errors"
	"unicode"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/google/uuid"
)

func (s *Service) GetUserReviews(
	ctx context.Context,
	userID uuid.UUID,
) ([]entities.PullRequest, error) {
	return nil, nil
}

func (s *Service) SetUserIsActive(
	ctx context.Context,
	userID []uuid.UUID,
	isActive bool,
) ([]entities.User, error) {
	return nil, nil
}

func validateUsername(username string) error {
	// длина: от 2 до 50 символов
	if len(username) < 2 {
		return errors.New("username must contain at least 2 characters")
	}
	if len(username) > 50 {
		return errors.New("username must contain at most 50 characters")
	}

	// только латинские буквы
	for _, r := range username {
		if !unicode.IsLetter(r) || r > unicode.MaxASCII {
			return errors.New("username can only contain english letters")
		}
	}

	return nil
}
