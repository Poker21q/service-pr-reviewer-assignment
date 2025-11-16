package entities

// Team представляет команду — группу пользователей с уникальным именем.
type Team struct {
	// Name уникальное имя команды.
	Name string

	// Members список пользователей команды.
	Members []User
}
