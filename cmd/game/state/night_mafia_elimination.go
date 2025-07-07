package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"strings"

	"github.com/teilomillet/gollm"
)

func (gs *GameState) nightMafiaEliminationVote(player *game.Player) error {
	prompt := gs.basePrompt(player)

	playerName, err := gs.llm.Generate(context.Background(), &prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response for mafia elimination vote: %w", err)
	}
	playerName = strings.TrimSpace(playerName)
	if playerName == "" {
		return fmt.Errorf("empty response received for mafia elimination vote")
	}

	gs.Conversation.AddMessage(
		player,
		playerName,
		enums.RoleMafia,
	)

	if cnt, exists := gs.votes[playerName]; exists {
		gs.votes[playerName] = cnt + 1
	} else {
		gs.votes[playerName] = 1
	}

	return nil
}

func (gs *GameState) nightMultipleMafiaElimination(mafias []*game.Player, nonMafiaPlayers []string) error {
	killMsg := fmt.Sprintf("Mafia members must vote to eliminate one of the following players: %s. The player with the most votes will be eliminated, if the votes are even, then a random player out of the proposed candidates will be eliminated. Respond with the player name and nothing else.", strings.Join(nonMafiaPlayers, ", "))

	gs.Conversation.AddMessage(
		game.NARRATOR,
		killMsg,
		enums.RoleMafia,
	)

	for _, mafia := range mafias {
		if err := gs.nightMafiaEliminationVote(mafia); err != nil {
			return fmt.Errorf("failed to process mafia elimination vote for %s: %w", mafia.Name, err)
		}
	}

	var maxVotes int
	for candidate, voteCount := range gs.votes {
		if voteCount > maxVotes {
			maxVotes = voteCount
			gs.mafiaElimination = candidate
		}
	}

	return nil
}

func (gs *GameState) nightSingleMafiaElimination(player *game.Player, nonMafiaPlayers []string) error {
	prompt := gs.basePrompt(player)

	killMsg := fmt.Sprintf("You are the only Mafia member left. Eliminate one of the following players: %s. Respond with the player name and nothing else.", strings.Join(nonMafiaPlayers, ", "))

	prompt.Messages = append(prompt.Messages, gollm.PromptMessage{
		Role:    "user",
		Name:    string(enums.RoleNarrator),
		Content: killMsg,
	})

	gs.Conversation.AddMessage(
		game.NARRATOR,
		killMsg,
		enums.RoleMafia,
	)

	playerName, err := gs.llm.Generate(context.Background(), &prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response for single mafia elimination: %w", err)
	}
	playerName = strings.TrimSpace(playerName)
	if playerName == "" {
		return fmt.Errorf("empty response received for single mafia elimination")
	}

	gs.Conversation.AddMessage(
		player,
		playerName,
		enums.RoleMafia,
	)

	if !gs.eliminatePlayer(playerName) {
		return fmt.Errorf("player %s does not exist", playerName)
	}

	return nil
}

func (gs *GameState) nightMafiaElimination() error {
	var mafiaPlayers []*game.Player

	for i := range gs.players {
		if gs.players[i].Role == enums.RoleMafia {
			mafiaPlayers = append(mafiaPlayers, &gs.players[i])
		}
	}

	nonMafiaPlayers := make([]string, 0)
	for i := range gs.players {
		if gs.players[i].Role != enums.RoleMafia {
			nonMafiaPlayers = append(nonMafiaPlayers, gs.players[i].Name)
		}
	}

	if len(mafiaPlayers) == 0 {
		return fmt.Errorf("no mafia players found during night elimination")
	} else if len(mafiaPlayers) == 1 {
		return gs.nightSingleMafiaElimination(mafiaPlayers[0], nonMafiaPlayers)
	} else {
		return gs.nightMultipleMafiaElimination(mafiaPlayers, nonMafiaPlayers)
	}
}
