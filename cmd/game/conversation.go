package game

import (
	"fmt"
	"mafia/cmd/enums"

	"github.com/fatih/color"
)

type ConversationLog struct {
	Player  *Player
	Message string
	Role    enums.Role // The message may be role-specific, for example night elimination phase for Mafia
}

type Conversation struct {
	conversation []ConversationLog
}

func (c *Conversation) AddMessage(player *Player, message string, role ...enums.Role) {
	c.conversation = append(c.conversation, ConversationLog{
		Player:  player,
		Message: message,
	})

	if len(role) > 0 {
		c.conversation[len(c.conversation)-1].Role = role[0]
	}

	args := []any{
		player.Name,
		message,
	}
	var format string
	if player.Role == enums.RoleNarrator {
		format = "%s: %s"
	} else {
		args = append([]any{player.Role}, args...)
		format = "(%s) %s: %s"
	}

	if len(role) > 0 {
		args = append([]any{role[0]}, args...)
		format = "**%s** " + format
	}

	var msg string
	switch player.Role {
	case enums.RoleNarrator:
		msg = color.MagentaString(format, args...)
	case enums.RoleMafia:
		msg = color.RedString(format, args...)
	case enums.RoleDoctor:
		msg = color.GreenString(format, args...)
	case enums.RoleCitizen:
		msg = color.CyanString(format, args...)
	case enums.RoleDetective:
		msg = color.BlueString(format, args...)
	}

	fmt.Println(msg)
}

func (c *Conversation) GetMessages() []ConversationLog {
	return c.conversation
}
