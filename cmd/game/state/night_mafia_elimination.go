package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/llm/models"
	"strings"
)

func (gs *GameState) nightMafiaEliminationVote(player *game.Player) error {
	messages := gs.baseMessages(player)

	response, err := gs.llm.Generate(context.Background(), messages)
	if err != nil {
		return fmt.Errorf("failed to generate response for mafia elimination vote: %w", err)
	}
	response.Content = strings.Trim(response.Content, " \n")
	if response.Content == "" {
		return fmt.Errorf("empty response received for mafia elimination vote")
	}

	gs.Conversation.AddMessage(
		player,
		response,
		enums.RoleMafia,
	)

	if cnt, exists := gs.votes[response.Content]; exists {
		gs.votes[response.Content] = cnt + 1
	} else {
		gs.votes[response.Content] = 1
	}

	return nil
}

func (gs *GameState) nightMultipleMafiaElimination(mafias []*game.Player) error {
	killMsg := "As mafia members, you and your partner must choose a peaceful player to eliminate tonight. Reply ONLY with the exact name of the player you wish to eliminate. Do not include any extra words or explanations."

	gs.Conversation.AddMessagePlaintext(
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
	messages := gs.baseMessages(player)

	killMsg := "As a mafia member, you must choose a peaceful player to eliminate tonight. Reply ONLY with the exact name of the player you wish to eliminate. Do not include any extra words or explanations."

	messages = append(messages, models.GenerateMessage{
		Role:    "user",
		Name:    string(enums.RoleNarrator),
		Content: killMsg,
	})

	gs.Conversation.AddMessagePlaintext(
		game.NARRATOR,
		killMsg,
		enums.RoleMafia,
	)

	response, err := gs.llm.Generate(context.Background(), messages)
	if err != nil {
		return fmt.Errorf("failed to generate response for single mafia elimination: %w", err)
	}
	response.Content = strings.Trim(response.Content, " \n")
	if response.Content == "" {
		return fmt.Errorf("empty response received for single mafia elimination")
	}

	gs.Conversation.AddMessage(
		player,
		response,
		enums.RoleMafia,
	)

	if !gs.eliminatePlayer(response.Content) {
		return fmt.Errorf("player %s does not exist", response.Content)
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
