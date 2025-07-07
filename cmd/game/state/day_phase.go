package state

import (
	"fmt"
	"mafia/cmd/game"
)

func (gs *GameState) DayPhase() error {
	gs.Conversation.AddMessage(
		game.NARRATOR,
		"The city is awake, the discussions are about to begin.",
	)

	for i := range gs.players {
		if err := gs.dayDiscussion(&gs.players[i]); err != nil {
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
			"The day ends with some accusations, we will start voting, cast your vote for a single person or abstain. The player with >50% of the votes will be eliminated. If no player has more than 50% of the votes, no one will be eliminated.",
		)
		if err := gs.dayVoting(); err != nil {
			return fmt.Errorf("failed to proceed with day voting: %w", err)
		}
	}

	gs.Clear()

	return nil
}
