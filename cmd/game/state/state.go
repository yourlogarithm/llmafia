package state

import (
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/llm"
	"mafia/cmd/llm/models"
)

type GameState struct {
	llm              llm.LLM
	players          []game.Player
	Cycle            int
	Conversation     game.Conversation
	accusedPlayers   map[string]string // map[accused]by_player
	votes            map[string]int    // either day elimination votes or night kill votes
	mafiaElimination string            // Will be killed at night by mafia if not saved by doc
	lastSaved        string
}

func NewGameState(players []game.Player, llm llm.LLM) *GameState {
	state := GameState{
		llm:            llm,
		players:        make([]game.Player, len(players)),
		accusedPlayers: make(map[string]string, len(players)),
		votes:          make(map[string]int, len(players)),
	}
	copy(state.players, players)
	state.Conversation.AddMessagePlaintext(
		game.NARRATOR,
		"The game has just started.",
	)
	return &state
}

func (gs *GameState) baseMessages(player *game.Player) (messages []models.GenerateMessage) {
	messages = append(messages, models.GenerateMessage{
		Role:    "user",
		Content: player.SystemPrompt,
	})

	for _, log := range gs.Conversation.GetMessages() {
		if log.Role == "" || log.Role == player.Role {
			messages = append(messages, models.GenerateMessage{
				Role:             "user",
				Name:             log.Player.Name,
				Content:          log.Message,
				ReasoningContent: log.Reasoning,
			})
		}
	}

	return messages
}

func (gs *GameState) EndgameStatus() enums.GameStatus {
	var mafiaCnt, peacefulCnt int
	for _, player := range gs.players {
		if player.Role == enums.RoleMafia {
			mafiaCnt++
		} else {
			peacefulCnt++
		}
	}
	if mafiaCnt == 0 {
		return enums.GameStatusPeacefulWin
	} else if mafiaCnt >= peacefulCnt {
		return enums.GameStatusMafiaWin
	}
	return enums.GameStatusOngoing
}

func (gs *GameState) eliminatePlayer(name string) bool {
	for i := 0; i < len(gs.players); i++ {
		if gs.players[i].Name == name {
			gs.players = append(gs.players[:i], gs.players[i+1:]...)
			return true
		}
	}
	return false
}

func (gs *GameState) Clear() {
	gs.accusedPlayers = make(map[string]string, len(gs.players))
	gs.votes = make(map[string]int, len(gs.players))
}

func (gs *GameState) UpdateCycle() {
	gs.Cycle = (gs.Cycle + 1) % len(gs.players)
}
