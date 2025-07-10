package game

import (
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/llm/models"
	"strings"

	"github.com/fatih/color"
)

type ConversationLog struct {
	Player    *Player    `json:"-"`
	Message   string     `json:"message"`
	Role      enums.Role `json:"role"`      // The message may be role-specific, for example night elimination phase for Mafia
	Reasoning string     `json:"reasoning"` // Optional reasoning content for the message
}

type Conversation struct {
	conversation []ConversationLog
}

func (c *Conversation) AddMessagePlaintext(player *Player, message string, role ...enums.Role) {
	response := models.GenerateResponse{
		Content:   message,
		Reasoning: "",
	}
	c.AddMessage(player, response, role...)
}

func (c *Conversation) AddMessage(player *Player, response models.GenerateResponse, role ...enums.Role) {
	response.Content = strings.Trim(response.Content, " \n")
	response.Reasoning = strings.Trim(response.Reasoning, " \n")
	c.conversation = append(c.conversation, ConversationLog{
		Player:    player,
		Message:   response.Content,
		Reasoning: response.Reasoning,
	})

	if len(role) > 0 {
		c.conversation[len(c.conversation)-1].Role = role[0]
	}

	args := []any{
		player.Name,
		response.Content,
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
