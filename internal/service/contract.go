package service

import (
	"context"

	"service-pr-reviewer-assignment/internal/service/entities"
)

type storage interface {
	CreateTeam(ctx context.Context, teamname string) (createdTeam *entities.Team, err error)

	CreateOrUpdateUsers(
		ctx context.Context,
		usersModify []entities.UserModify) (updatedUsers []entities.User, err error)
}

type txManager interface {
	Read(ctx context.Context, fn func(ctx context.Context) error) error
	Write(ctx context.Context, fn func(ctx context.Context) error) error
}
