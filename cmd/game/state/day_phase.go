package state

import (
	"fmt"
	"mafia/cmd/game"
	"strings"
)

func rotatePlayers(s []game.Player, i int) []game.Player {
	i %= len(s)
	n := len(s)
	if n == 0 {
		return s
	}
	i = ((i % n) + n) % n
	return append(s[i:], s[:i]...)
}

func (gs *GameState) DayPhase(firstDay bool) error {
	var message string

	gs.players = rotatePlayers(gs.players, gs.cycle)

	var playerOrder []string
	for _, player := range gs.players {
		playerOrder = append(playerOrder, player.Name)
	}
	playerOrderStr := fmt.Sprintf("Players will speak uninterrupted one by one in the following order: %s", strings.Join(playerOrder, ", "))

	if firstDay {
		message = fmt.Sprintf("The first day has begun. The discussions are about to start. %s", playerOrderStr)
	} else {
		message = fmt.Sprintf("The day has begun. The discussions are about to start. %s", playerOrderStr)
	}
	gs.Conversation.AddMessagePlaintext(
		game.NARRATOR,
		message,
	)

	for i := 0; i < len(gs.players); i++ {
		if err := gs.dayDiscussion(&gs.players[i], i == 0); err != nil {
			return err
		}
	}

	if len(gs.accusedPlayers) == 0 {
		gs.Conversation.AddMessagePlaintext(
			game.NARRATOR,
			"No player has been accused of being a Mafia member, so the day ends without any elimination voting.",
		)
	} else {
		var accusedNames []string
		for accusedName := range gs.accusedPlayers {
			accusedNames = append(accusedNames, accusedName)
		}
		gs.Conversation.AddMessagePlaintext(
			game.NARRATOR,
			fmt.Sprintf("Players: %s - have been accused. It's time to vote, each player may vote to eliminate one accused person or abstain. If any accused receives more than 50%% of the votes, they will be eliminated. Otherwise, no one will be eliminated this day.", strings.Join(accusedNames, ", ")),
		)
		if err := gs.dayVoting(); err != nil {
			return fmt.Errorf("failed to proceed with day voting: %w", err)
		}
	}

	gs.Clear()

	return nil
}
