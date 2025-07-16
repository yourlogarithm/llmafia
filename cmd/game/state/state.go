package state

import (
	"encoding/json"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/llm"
	"mafia/cmd/llm/models"
	"os"
	"time"
)

type GameState struct {
	llm              llm.LLM
	players          []game.Player
	cycle            int
	accusedPlayers   map[string]string // map[accused]by_player
	votes            map[string]int    // either day elimination votes or night kill votes
	mafiaElimination string            // Will be killed at night by mafia if not saved by doc
	lastSaved        string
	Conversation     game.Conversation
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
	gs.cycle++
}

func (gs *GameState) Dump(t time.Time, players []game.Player, f *os.File) error {
	var conversation []ConversationLogPlayerName
	for _, msg := range gs.Conversation.GetMessages() {
		conversation = append(conversation, ConversationLogPlayerName{
			Player:    msg.Player.Name,
			Message:   msg.Message,
			Role:      msg.Role,
			Reasoning: msg.Reasoning,
		})
	}
	game_log := GameLog{
		Timestamp:    t,
		Players:      players,
		Conversation: conversation,
		Cycles:       gs.cycle,
		Status:       gs.EndgameStatus(),
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(game_log)
}
