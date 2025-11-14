package storage

import (
	"context"

	"service-pr-reviewer-assignment/internal/service/entities"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
)

func (s *Storage) CreateTeam(ctx context.Context, teamname string) (*entities.Team, error) {
	const query = `INSERT INTO teams (name) VALUES ($1)`

	row := s.querier.QueryRow(ctx, query, teamname)

	var res struct {
		Name string
	}

	if err := row.Scan(&res.Name); err != nil {
		return nil, err
	}

	return &entities.Team{
		Name: res.Name,
	}, nil
}

func (s *Storage) CreateOrUpdateUsers(ctx context.Context, users []entities.UserModify) ([]entities.User, error) {
	if len(users) == 0 {
		return nil, nil
	}

	builder := sq.Insert("users").
		Columns("id", "name", "team_name", "is_active").
		PlaceholderFormat(sq.Dollar)

	for _, u := range users {
		id := u.ID
		if id == nil {
			newID := uuid.New()
			id = &newID
		}
		builder = builder.Values(*id, u.Name, u.TeamName, u.IsActive)
	}

	query, args, err := builder.Suffix(`
        ON CONFLICT (id) DO UPDATE
        SET
            name = EXCLUDED.name,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
        RETURNING id, name, team_name, is_active
    `).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.querier.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.User
	for rows.Next() {
		var m struct {
			ID       uuid.UUID
			Name     string
			TeamName string
			IsActive bool
		}
		err = rows.Scan(&m.ID, &m.Name, &m.TeamName, &m.IsActive)
		if err != nil {
			return nil, err
		}

		result = append(result, entities.User{
			ID:       m.ID,
			Name:     m.Name,
			TeamName: m.TeamName,
			IsActive: m.IsActive,
		})
	}

	return result, nil
}
