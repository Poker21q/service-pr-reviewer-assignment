package entities

import "github.com/google/uuid"

// Team представляет команду — группу пользователей с уникальным именем.
type Team struct {
	// Name уникальное имя команды.
	Name string

	// MembersIDs список пользователей команды.
	MembersIDs []uuid.UUID
}
