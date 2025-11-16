package converters

import (
	"service-pr-reviewer-assignment/internal/generated/api/dto"
	"service-pr-reviewer-assignment/internal/service/entities"
)

func UserToDTO(user *entities.User) *dto.User {
	return &dto.User{
		UserId:   user.ID,
		Username: user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
