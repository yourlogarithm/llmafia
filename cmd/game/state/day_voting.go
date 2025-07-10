package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/llm/models"
	"mafia/cmd/utils"
	"strings"
	"sync"
)

func (gs *GameState) collectPlayerVote(player *game.Player) error {
	messages := gs.baseMessages(player)

	var accusationStrings []string
	for accused, byPlayer := range gs.accusedPlayers {
		accusationStrings = append(accusationStrings, fmt.Sprintf("- %s, accused by %s", accused, byPlayer))
	}

	messages = append(messages, models.GenerateMessage{
		Role:    "user",
		Name:    string(enums.RoleNarrator),
		Content: "You must choose to vote to eliminate an accused player on the basis that he/she is a Mafia member, or you can abstain. The elimination candidates are as follows:\n" + strings.Join(accusationStrings, "\n") + "\nRespond in the following JSON format: {\"vote\": \"candidate_name\"} or {\"abstain\": true} if you choose not to vote for anybody.",
	})

	var voteResponse struct {
		Vote    string `json:"vote,omitempty"`
		Abstain bool   `json:"abstain,omitempty"`
	}

	response, err := gs.llm.Generate(context.Background(), messages)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	if response.Content == "" {
		response.Content = "I abstain from voting."
		gs.Conversation.AddMessage(player, response)
		return nil
	} else if err := utils.ParseJSONResponsePermissive(response.Content, &voteResponse); err != nil {
		return fmt.Errorf("failed to parse vote response: %w", err)
	}

	if voteResponse.Abstain {
		response.Content = "I abstain from voting."
		gs.Conversation.AddMessage(player, response)
	} else {
		candidate := strings.Trim(voteResponse.Vote, " \n")
		if candidate == "" {
			return fmt.Errorf("invalid vote: empty candidate name")
		}

		if _, exists := gs.accusedPlayers[candidate]; !exists {
			return fmt.Errorf("invalid vote: %s is not an accused player", candidate)
		}
		response.Content = fmt.Sprintf("I vote to eliminate %s.", candidate)
		gs.Conversation.AddMessage(player, response)
		if cnt, exists := gs.votes[candidate]; exists {
			gs.votes[candidate] = cnt + 1
		} else {
			gs.votes[candidate] = 1
		}
	}

	return nil
}

func (gs *GameState) dayVoting() (err error) {
	var wg sync.WaitGroup

	for i := range gs.players {
		wg.Add(1)
		go func(p *game.Player) {
			defer wg.Done()
			localErr := gs.collectPlayerVote(p)
			if localErr != nil {
				err = localErr
			}
		}(&gs.players[i])
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
		gs.Conversation.AddMessagePlaintext(game.NARRATOR, fmt.Sprintf("%s has been eliminated with %d votes.", eliminatedPlayer, maxVotes))
		if !gs.eliminatePlayer(eliminatedPlayer) {
			return fmt.Errorf("player %s does not exist", eliminatedPlayer)
		}
	} else {
		gs.Conversation.AddMessagePlaintext(game.NARRATOR, fmt.Sprintf("The highest vote count was %d, which is not enough to eliminate anyone, thus the day ends without any elimination.", maxVotes))
	}

	return nil
}
