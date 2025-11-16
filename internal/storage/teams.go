package storage

import (
	"context"
	"errors"
	"fmt"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/jackc/pgx/v5/pgconn"
)

type teamDB struct {
	Name string
}

func (s *Storage) CreateTeam(ctx context.Context, teamName string) (*entities.Team, error) {
	const query = `INSERT INTO teams (name) VALUES ($1) RETURNING name`

	var teamDB teamDB
	err := s.querier.QueryRow(ctx, query, teamName).Scan(&teamDB.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return nil, &entities.ErrTeamAlreadyExists{Name: teamName}
		}
		return nil, fmt.Errorf("create team: %w", err)
	}

	team := convertTeamDBToEntity(teamDB)
	return &team, nil
}

func (s *Storage) IsTeamExists(ctx context.Context, teamName string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`

	var exists bool
	err := s.querier.QueryRow(ctx, query, teamName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check team exists: %w", err)
	}

	return exists, nil
}

func convertTeamDBToEntity(teamDB teamDB) entities.Team {
	return entities.Team{
		Name: teamDB.Name,
	}
}
