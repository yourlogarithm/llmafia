package state

import (
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"time"
)

type GameLog struct {
	Timestamp    time.Time                   `json:"timestamp"`
	Players      []game.Player               `json:"players"`
	Conversation []ConversationLogPlayerName `json:"conversation"`
	Cycles       int                         `json:"cycles"`
	Status       enums.GameStatus            `json:"status"`
}

type ConversationLogPlayerName struct {
	Player    string     `json:"player"`
	Message   string     `json:"message"`
	Role      enums.Role `json:"role"`      // The message may be role-specific, for example night elimination phase for Mafia
	Reasoning string     `json:"reasoning"` // Optional reasoning content for the message
}
