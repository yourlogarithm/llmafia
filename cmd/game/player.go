package game

import "mafia/cmd/enums"

type Player struct {
	Name         string     `json:"name"`
	Role         enums.Role `json:"role"`
	SystemPrompt string     `json:"-"`
}

var NARRATOR = &Player{
	Name: string(enums.RoleNarrator),
	Role: enums.RoleNarrator,
}
