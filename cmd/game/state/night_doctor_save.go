package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"strings"
)

func (gs *GameState) nightDoctorSave() error {
	var doctor *game.Player
	for i := range gs.players {
		p := &gs.players[i]
		if p.Role == enums.RoleDoctor {
			doctor = p
			break
		}
	}
	if doctor == nil {
		return nil
	}

	var rowMessage string
	if gs.lastSaved != "" {
		rowMessage = fmt.Sprintf(" Remember that you cannot save %s two times in a row.", gs.lastSaved)
	}

	gs.Conversation.AddMessage(
		game.NARRATOR,
		fmt.Sprintf("The doctor must choose someone to protect from elimination tonight.%s Choose by responding with the name of the player you want to save and nothing else.", rowMessage),
		enums.RoleDoctor,
	)

	prompt := gs.basePrompt(doctor)
	playerName, err := gs.llm.Generate(context.Background(), &prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response for doctor save: %w", err)
	}

	playerName = strings.TrimSpace(playerName)
	if playerName == "" {
		return fmt.Errorf("empty response received for doctor save")
	}

	gs.Conversation.AddMessage(
		doctor,
		playerName,
		enums.RoleDoctor,
	)

	gs.lastSaved = playerName

	if gs.mafiaElimination == playerName {
		gs.Conversation.AddMessage(
			game.NARRATOR,
			"Tonight's victim has been saved by the doctor.",
		)
		gs.mafiaElimination = ""
	}

	return nil
}
