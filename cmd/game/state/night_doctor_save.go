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
	if gs.lastSaved == doctor.Name {
		rowMessage = " Remember that you cannot save yourself two times in a row."
	} else if gs.lastSaved != "" {
		rowMessage = fmt.Sprintf(" Remember that you cannot save %s two times in a row.", gs.lastSaved)
	}

	gs.Conversation.AddMessagePlaintext(
		game.NARRATOR,
		fmt.Sprintf("As the doctor, you must choose someone to protect from elimination tonight.%s Reply ONLY with the exact name of the player you wish to protect. Do not include any extra words or explanations.", rowMessage),
		enums.RoleDoctor,
	)

	messages := gs.baseMessages(doctor)
	response, err := gs.llm.Generate(context.Background(), messages)
	if err != nil {
		return fmt.Errorf("failed to generate response for doctor save: %w", err)
	}

	response.Content = strings.Trim(response.Content, " \n")
	if response.Content == "" {
		return fmt.Errorf("empty response received for doctor save")
	}

	if gs.lastSaved != "" && response.Content == gs.lastSaved {
		var message string
		if gs.lastSaved == doctor.Name {
			message = "You cannot save yourself again tonight."
		} else {
			message = fmt.Sprintf("You cannot save %s again tonight.", gs.lastSaved)
		}
		gs.Conversation.AddMessagePlaintext(
			game.NARRATOR,
			message+" Because you forgot that - no one will be saved tonight.",
		)
		gs.lastSaved = ""
		return nil
	}

	gs.Conversation.AddMessage(
		doctor,
		response,
		enums.RoleDoctor,
	)

	gs.lastSaved = response.Content

	if gs.mafiaElimination == response.Content {
		gs.Conversation.AddMessagePlaintext(
			game.NARRATOR,
			"Tonight's victim has been saved by the doctor.",
		)
		gs.mafiaElimination = ""
	}

	return nil
}
