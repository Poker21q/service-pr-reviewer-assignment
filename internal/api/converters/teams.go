package converters

import (
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/service/entities"
)

func TeamFromDTO(req dto.Team) (string, []entities.User, error) {
	teamName := req.TeamName

	users := make([]entities.User, 0, len(req.Members))
	for _, m := range req.Members {
		uid := m.UserId
		name := m.Username
		team := teamName
		isActive := m.IsActive

		users = append(users, entities.User{
			ID:       uid,
			Name:     name,
			TeamName: team,
			IsActive: isActive,
		})
	}

	return teamName, users, nil
}

func TeamToDTO(team *entities.Team) dto.Team {
	members := make([]dto.TeamMember, 0, len(team.Members))
	for _, u := range team.Members {
		members = append(members, dto.TeamMember{
			UserId:   u.ID,
			Username: u.Name,
			IsActive: u.IsActive,
		})
	}

	return dto.Team{
		TeamName: team.Name,
		Members:  members,
	}
}
