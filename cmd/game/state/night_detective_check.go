package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"strings"
)

func (gs *GameState) nightDetectiveCheck() error {
	var detective *game.Player
	for i := range gs.players {
		p := &gs.players[i]
		if p.Role == enums.RoleDetective {
			detective = p
			break
		}
	}
	if detective == nil {
		return nil
	}

	gs.Conversation.AddMessagePlaintext(
		game.NARRATOR,
		"As the detective, you must investigate one player tonight. Reply ONLY with the exact name of the player you wish to investigate. Do not include any extra words or explanations.",
		enums.RoleDetective,
	)
	messages := gs.baseMessages(detective)
	response, err := gs.llm.Generate(context.Background(), messages)
	if err != nil {
		return fmt.Errorf("failed to generate response for detective check: %w", err)
	}

	response.Content = strings.Trim(response.Content, " \n")
	if response.Content == "" {
		return fmt.Errorf("empty response received for detective check")
	}

	gs.Conversation.AddMessage(
		detective,
		response,
		enums.RoleDetective,
	)

	var found bool
	for i := range gs.players {
		p := &gs.players[i]
		if p.Name == response.Content {
			if p.Role == enums.RoleMafia {
				gs.Conversation.AddMessagePlaintext(
					game.NARRATOR,
					fmt.Sprintf("%s is a Mafia member", response.Content),
					enums.RoleDetective,
				)
			} else {
				gs.Conversation.AddMessagePlaintext(
					game.NARRATOR,
					fmt.Sprintf("%s is NOT a Mafia member.", response.Content),
					enums.RoleDetective,
				)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("player %s not found in the game state", response.Content)
	}

	return nil
}
