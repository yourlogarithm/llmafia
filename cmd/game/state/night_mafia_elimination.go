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

func (gs *GameState) nightMultipleMafiaElimination(mafias []*game.Player) error {
	killMsg := "As mafia members, you and your partner must choose a peaceful player to eliminate tonight. Reply ONLY with the exact name of the player you wish to eliminate. Do not include any extra words or explanations."

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

func (gs *GameState) nightSingleMafiaElimination(player *game.Player) error {
	prompt := gs.basePrompt(player)

	killMsg := "As a mafia member, you must choose a peaceful player to eliminate tonight. Reply ONLY with the exact name of the player you wish to eliminate. Do not include any extra words or explanations."

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

	if len(mafiaPlayers) == 0 {
		return fmt.Errorf("no mafia players found during night elimination")
	} else if len(mafiaPlayers) == 1 {
		return gs.nightSingleMafiaElimination(mafiaPlayers[0])
	} else {
		return gs.nightMultipleMafiaElimination(mafiaPlayers)
	}
}
