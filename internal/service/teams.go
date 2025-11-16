package service

import (
	"context"
	"fmt"
	"unicode"

	"service-pr-reviewer-assignment/internal/service/entities"
)

func (s *Service) CreateTeam(
	ctx context.Context,
	teamName string,
	members []entities.User,
) (*entities.Team, error) {
	if err := validateTeamName(teamName); err != nil {
		return nil, fmt.Errorf("validate team name: %w", err)
	}
	if err := validateMembers(members); err != nil {
		return nil, fmt.Errorf("validate members: %w", err)
	}
	if err := checkDuplicateUserIDs(members); err != nil {
		return nil, fmt.Errorf("check duplicate users: %w", err)
	}

	var updatedMembers []entities.User
	err := s.txManager.Write(ctx, func(ctx context.Context) error {
		if _, err := s.storage.CreateTeam(ctx, teamName); err != nil {
			return fmt.Errorf("create team: %w", err)
		}

		for i := range members {
			members[i].TeamName = teamName
		}

		var err error
		updatedMembers, err = s.storage.CreateOrUpdateUsers(ctx, members)
		if err != nil {
			return fmt.Errorf("create or update users: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("create team transaction: %w", err)
	}

	return &entities.Team{
		Name:    teamName,
		Members: updatedMembers,
	}, nil
}

func (s *Service) GetTeam(
	ctx context.Context,
	teamName string,
) (*entities.Team, error) {
	var members []entities.User
	err := s.txManager.Read(ctx, func(ctx context.Context) error {
		exists, err := s.storage.IsTeamExists(ctx, teamName)
		if err != nil {
			return fmt.Errorf("check team exists: %w", err)
		}
		if !exists {
			return &entities.ErrTeamNotFound{Name: teamName}
		}

		members, err = s.storage.GetUsersByTeamName(ctx, teamName)
		if err != nil {
			return fmt.Errorf("get team members: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	return &entities.Team{
		Name:    teamName,
		Members: members,
	}, nil
}

func validateMembers(members []entities.User) error {
	for _, member := range members {
		if err := validateUsername(member.Name); err != nil {
			return &entities.ErrUserNameValidation{Reason: member.Name}
		}
	}
	return nil
}

func validateTeamName(teamName string) error {
	if len(teamName) < 2 {
		return &entities.ErrTeamNameValidation{Reason: "team name too short"}
	}
	if len(teamName) > 50 {
		return &entities.ErrTeamNameValidation{Reason: "team name too long"}
	}

	for _, r := range teamName {
		if !unicode.IsLetter(r) || r > unicode.MaxASCII {
			return &entities.ErrTeamNameValidation{Reason: "team name must contain only english letters"}
		}
	}

	return nil
}
