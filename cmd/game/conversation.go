package game

import (
	"mafia/cmd/enums"

	"github.com/fatih/color"
)

type ConversationLog struct {
	Player  Player
	Message string
	Role    enums.Role // The message may be role-specific, for example night elimination phase for Mafia
}

type Conversation struct {
	conversation []ConversationLog
}

func (c *Conversation) AddMessage(player Player, message string) {
	c.conversation = append(c.conversation, ConversationLog{
		Player:  player,
		Message: message,
	})
	switch player.Role {
	case enums.RoleNarrator:
		color.Magenta("%s (%s): %s\n", player.Name, player.Role, message)
	case enums.RoleMafia:
		color.Red("%s (%s): %s\n", player.Name, player.Role, message)
	case enums.RoleDoctor:
		color.Green("%s (%s): %s\n", player.Name, player.Role, message)
	case enums.RoleCitizen:
		color.Cyan("%s (%s): %s\n", player.Name, player.Role, message)
	case enums.RoleDetective:
		color.Blue("%s (%s): %s\n", player.Name, player.Role, message)
	}
}

func (c *Conversation) GetMessages() []ConversationLog {
	return c.conversation
}
