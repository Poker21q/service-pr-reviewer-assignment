package converters

import (
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/google/uuid"
)

func TeamFromDTO(req dto.Team) (string, []entities.UserModify, error) {
	teamName := req.TeamName

	users := make([]entities.UserModify, 0, len(req.Members))
	for _, m := range req.Members {
		uid, err := uuid.Parse(m.UserId)
		if err != nil {
			return "", nil, err
		}

		name := m.Username
		team := teamName
		isActive := m.IsActive

		users = append(users, entities.UserModify{
			ID:       &uid,
			Name:     &name,
			TeamName: &team,
			IsActive: &isActive,
		})
	}

	return teamName, users, nil
}

func TeamToDTO(teamName string, teamMembers []entities.User) dto.Team {
	members := make([]dto.TeamMember, 0, len(teamMembers))
	for _, u := range teamMembers {
		members = append(members, dto.TeamMember{
			UserId:   u.ID.String(),
			Username: u.Name,
			IsActive: u.IsActive,
		})
	}

	return dto.Team{
		TeamName: teamName,
		Members:  members,
	}
}
