package storage

import (
	"context"
	"errors"
	"fmt"

	"service-pr-reviewer-assignment/internal/service/entities"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userDB struct {
	ID       uuid.UUID
	Name     string
	TeamName string
	IsActive bool
}

func (s *Storage) CreateOrUpdateUsers(ctx context.Context, users []entities.User) ([]entities.User, error) {
	if len(users) == 0 {
		return nil, nil
	}

	builder := s.stmtBuilder.
		Insert("users").
		Columns("id", "name", "team_name", "is_active")

	for _, user := range users {
		if user.ID == uuid.Nil {
			user.ID = uuid.New()
		}
		builder = builder.Values(user.ID, user.Name, user.TeamName, user.IsActive)
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
		return nil, fmt.Errorf("build upsert query: %w", err)
	}

	rows, err := s.querier.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("upsert users: %w", err)
	}
	defer rows.Close()

	var usersDB []userDB
	for rows.Next() {
		var userDB userDB
		if err := rows.Scan(&userDB.ID, &userDB.Name, &userDB.TeamName, &userDB.IsActive); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		usersDB = append(usersDB, userDB)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return convertUsersDBToEntities(usersDB), nil
}

func (s *Storage) GetUsersByTeamName(ctx context.Context, teamName string) ([]entities.User, error) {
	const query = `SELECT id, name, team_name, is_active FROM users WHERE team_name = $1`

	rows, err := s.querier.Query(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("query users by team: %w", err)
	}
	defer rows.Close()

	var usersDB []userDB
	for rows.Next() {
		var userDB userDB
		if err := rows.Scan(&userDB.ID, &userDB.Name, &userDB.TeamName, &userDB.IsActive); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		usersDB = append(usersDB, userDB)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return convertUsersDBToEntities(usersDB), nil
}

func (s *Storage) IsUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := s.querier.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user exists: %w", err)
	}

	return exists, nil
}

func (s *Storage) GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	const query = `SELECT id, name, team_name, is_active FROM users WHERE id = $1`

	var userDB userDB
	err := s.querier.QueryRow(ctx, query, userID).Scan(&userDB.ID, &userDB.Name, &userDB.TeamName, &userDB.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &entities.ErrUserNotFound{UserID: pointer.To(userID)}
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	user := convertUserDBToEntity(userDB)
	return &user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	const query = `
        UPDATE users 
        SET name = $2, team_name = $3, is_active = $4
        WHERE id = $1
        RETURNING id, name, team_name, is_active
    `

	var userDB userDB
	err := s.querier.QueryRow(
		ctx,
		query,
		user.ID,
		user.Name,
		user.TeamName,
		user.IsActive,
	).Scan(
		&userDB.ID,
		&userDB.Name,
		&userDB.TeamName,
		&userDB.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	updatedUser := convertUserDBToEntity(userDB)
	return &updatedUser, nil
}

func convertUsersDBToEntities(usersDB []userDB) []entities.User {
	out := make([]entities.User, 0, len(usersDB))
	for _, userDB := range usersDB {
		out = append(out, convertUserDBToEntity(userDB))
	}
	return out
}

func convertUserDBToEntity(userDB userDB) entities.User {
	return entities.User{
		ID:       userDB.ID,
		Name:     userDB.Name,
		TeamName: userDB.TeamName,
		IsActive: userDB.IsActive,
	}
}
