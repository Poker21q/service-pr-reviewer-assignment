package entities

import "github.com/google/uuid"

// User представляет участника команды.
type User struct {
	// ID уникальный идентификатор пользователя.
	ID uuid.UUID

	// Name имя пользователя.
	Name string

	// TeamName идентификатор команды пользователя.
	TeamName string

	// IsActive флаг активности пользователя.
	IsActive bool
}
