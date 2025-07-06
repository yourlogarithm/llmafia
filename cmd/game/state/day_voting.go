package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/utils"
	"strings"
	"sync"

	"github.com/teilomillet/gollm"
)

func (gs *GameState) CollectPlayerVote(player game.Player) error {
	prompt := gs.BasePrompt(player)

	var accusationStrings []string
	for accused, byPlayer := range gs.accusedPlayers {
		accusationStrings = append(accusationStrings, fmt.Sprintf("- %s, accused by %s", accused, byPlayer))
	}

	prompt.Messages = append(prompt.Messages, gollm.PromptMessage{
		Role:    "user",
		Name:    string(enums.RoleNarrator),
		Content: "You must choose to vote to eliminate an accused player on the basis that he/she is a Mafia member, or you can abstain. The elimination candidates are as follows:\n" + strings.Join(accusationStrings, "\n") + "\nRespond in the following JSON format: {\"vote\": \"candidate_name\"} or {\"abstain\": true} if you choose not to vote for anybody.",
	})

	var voteResponse struct {
		Vote    string `json:"vote,omitempty"`
		Abstain bool   `json:"abstain,omitempty"`
	}

	response, err := gs.llm.Generate(context.Background(), &prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	if response == "" {
		gs.Conversation.AddMessage(player, "I abstain from voting.")
		return nil
	} else if err := utils.ParseJSONResponsePermissive(response, &voteResponse); err != nil {
		return fmt.Errorf("failed to parse vote response: %w", err)
	}

	if voteResponse.Abstain {
		gs.Conversation.AddMessage(player, "I abstain from voting.")
	} else {
		candidate := strings.TrimSpace(voteResponse.Vote)
		if candidate == "" {
			return fmt.Errorf("invalid vote: empty candidate name")
		}

		if _, exists := gs.accusedPlayers[candidate]; !exists {
			return fmt.Errorf("invalid vote: %s is not an accused player", candidate)
		}

		gs.Conversation.AddMessage(player, fmt.Sprintf("I vote to eliminate %s.", candidate))
		if cnt, exists := gs.votes[candidate]; exists {
			gs.votes[candidate] = cnt + 1
		} else {
			gs.votes[candidate] = 1
		}
	}

	return nil
}

func (gs *GameState) DayVoting() error {
	gs.phase = enums.PhaseDayVoting

	var err error

	var wg sync.WaitGroup

	for _, player := range gs.players {
		wg.Add(1)
		go func(p game.Player) {
			defer wg.Done()
			err = gs.CollectPlayerVote(p)
		}(player)
	}

	wg.Wait()

	if err != nil {
		return fmt.Errorf("error collecting votes: %w", err)
	}

	var maxVotes int
	var eliminatedPlayer string
	for candidate, voteCount := range gs.votes {
		if voteCount > maxVotes {
			maxVotes = voteCount
			eliminatedPlayer = candidate
		}
	}

	prc := float64(maxVotes) / float64(len(gs.players))
	if prc > 0.5 {
		gs.Conversation.AddMessage(game.NARRATOR, fmt.Sprintf("%s has been eliminated with %d votes (%.2f%% of the total votes).", eliminatedPlayer, maxVotes, prc*100))
		for i, p := range gs.players {
			if p.Name == eliminatedPlayer {
				gs.players = append(gs.players[:i], gs.players[i+1:]...)
				break
			}
		}
	} else {
		gs.Conversation.AddMessage(game.NARRATOR, fmt.Sprintf("No player has been eliminated. The highest vote count was %d (%.2f%% of the total votes), which is not enough to eliminate anyone.", maxVotes, prc*100))
	}

	gs.accusedPlayers = make(map[string]string, len(gs.players))
	gs.votes = make(map[string]int, len(gs.players))

	return nil
}
