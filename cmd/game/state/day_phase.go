package state

import (
	"fmt"
	"mafia/cmd/game"
)

func (gs *GameState) DayPhase(firstDay bool) error {
	var message string
	if firstDay {
		message = "The first day has begun. the discussion is about to start."
	} else {
		message = "The next day has begun. the discussion is about to start."
	}
	gs.Conversation.AddMessage(
		game.NARRATOR,
		message,
	)

	for i := range gs.players {
		if err := gs.dayDiscussion(&gs.players[i], i == gs.Cycle); err != nil {
			return err
		}
	}

	if len(gs.accusedPlayers) == 0 {
		gs.Conversation.AddMessage(
			game.NARRATOR,
			"No player has been accused of being a Mafia member, so the day ends without any elimination voting.",
		)
	} else {
		gs.Conversation.AddMessage(
			game.NARRATOR,
			"Accusations have been made. It's time to vote: each player may vote to eliminate one accused person or abstain. If any accused receives more than 50% of the votes, they will be eliminated. Otherwise, no one will be eliminated this day.",
		)
		if err := gs.dayVoting(); err != nil {
			return fmt.Errorf("failed to proceed with day voting: %w", err)
		}
	}

	gs.Clear()

	return nil
}
