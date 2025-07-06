package state

import (
	"mafia/cmd/enums"
	"mafia/cmd/game"

	"github.com/teilomillet/gollm"
)

type GameState struct {
	llm            gollm.LLM
	players        []game.Player
	phase          enums.Phase
	Conversation   game.Conversation
	accusedPlayers map[string]string // map[accused]by_player
	votes          map[string]int    // either day elimination votes or night kill votes
}

func NewGameState(players []game.Player, llm gollm.LLM) *GameState {
	state := GameState{
		llm:            llm,
		players:        make([]game.Player, len(players)),
		accusedPlayers: make(map[string]string, len(players)),
		votes:          make(map[string]int, len(players)),
	}
	copy(state.players, players)
	state.Conversation.AddMessage(
		game.NARRATOR,
		"The game has just started. The city is awake, the discussions are about to begin.",
	)
	return &state
}

func (gs *GameState) BasePrompt(player game.Player) gollm.Prompt {
	prompt := gollm.Prompt{
		Messages: []gollm.PromptMessage{
			{
				Role:    "system",
				Content: player.SystemPrompt,
			},
		},
	}

	for _, log := range gs.Conversation.GetMessages() {
		if log.Role == "" || log.Role == player.Role {
			prompt.Messages = append(prompt.Messages, gollm.PromptMessage{
				Role:    "user",
				Name:    log.Player.Name,
				Content: log.Message,
			})
		}
	}

	return prompt
}
