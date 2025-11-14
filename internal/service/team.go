package service

import (
	"context"
	"errors"
	"unicode"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/AlekSi/pointer"
)

func (s *Service) CreateTeamByUsers(
	ctx context.Context,
	teamName string,
	members []entities.UserModify,
) ([]entities.User, error) {
	if err := validateTeamName(teamName); err != nil {
		return nil, err
	}

	for _, m := range members {
		username := pointer.Get(m.Name)
		if err := validateUsername(username); err != nil {
			return nil, err
		}
	}

	var updatedMembers []entities.User
	err := s.txManager.Write(ctx, func(ctx context.Context) error {
		_, err := s.storage.CreateTeam(ctx, teamName)
		if err != nil {
			return err
		}

		// Привязываем пользователей к команде
		for i := range members {
			members[i].TeamName = pointer.To(teamName)
		}

		updatedMembers, err = s.storage.CreateOrUpdateUsers(ctx, members)
		return err
	})
	if err != nil {
		return nil, err
	}

	return updatedMembers, nil
}

func validateTeamName(teamName string) error {
	// длина: от 2 до 50 символов
	if len(teamName) < 2 {
		return errors.New("teamname must contain at least 2 characters")
	}
	if len(teamName) > 50 {
		return errors.New("teamname must contain at most 50 characters")
	}

	// только латинские буквы
	for _, r := range teamName {
		if !unicode.IsLetter(r) || r > unicode.MaxASCII {
			return errors.New("teamname can only contain english letters")
		}
	}

	return nil
}
