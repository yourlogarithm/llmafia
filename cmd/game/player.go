package game

import "mafia/cmd/enums"

type Player struct {
	Name         string
	Role         enums.Role
	SystemPrompt string
}

var NARRATOR = Player{
	Name: string(enums.RoleNarrator),
	Role: enums.RoleNarrator,
}
