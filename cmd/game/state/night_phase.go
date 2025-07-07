package state

import (
	"fmt"
	"mafia/cmd/game"
)

func (gs *GameState) NightPhase() error {
	gs.Conversation.AddMessage(
		game.NARRATOR,
		"The night falls, the city is asleep, and the Mafia is about to make their move.",
	)

	err := gs.nightMafiaElimination()
	if err != nil {
		return fmt.Errorf("failed to proceed with night mafia elimination: %w", err)
	}

	err = gs.nightDoctorSave()
	if err != nil {
		return fmt.Errorf("failed to proceed with night doctor save: %w", err)
	}

	err = gs.nightDetectiveCheck()
	if err != nil {
		return fmt.Errorf("failed to proceed with night detective check: %w", err)
	}

	if gs.mafiaElimination != "" {
		if !gs.eliminatePlayer(gs.mafiaElimination) {
			return fmt.Errorf("failed to eliminate player %s", gs.mafiaElimination)
		}
		gs.Conversation.AddMessage(
			game.NARRATOR,
			fmt.Sprintf("Tonight %s has been eliminated.", gs.mafiaElimination),
		)
	}

	gs.Clear()

	return nil
}
