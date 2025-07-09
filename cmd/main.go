package main

import (
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/game/state"
	"mafia/cmd/llm"
	"mafia/cmd/prompts"
	"strings"
)

func main() {
	names := []string{
		"Alice", "Bob", "Charlie", "Diana", "Ethan", "Fiona", "George",
	}
	// names = names[:2]
	// rand.Shuffle(len(names), func(i, j int) {
	// 	names[i], names[j] = names[j], names[i]
	// })

	players := make([]game.Player, len(names))
	for i, name := range names {
		restPlayers := make([]string, 0, len(names)-1)
		for j, otherName := range names {
			if j != i {
				restPlayers = append(restPlayers, otherName)
			}
		}
		restPlayersStr := strings.Join(restPlayers, ", ")

		if i < 2 {
			players[i].Role = enums.RoleMafia
			players[i].SystemPrompt = fmt.Sprintf(prompts.MAFIA, name, restPlayersStr, enums.RoleMafia, names[(i+1)%2])
		} else if i == 2 {
			players[i].Role = enums.RoleDoctor
			players[i].SystemPrompt = fmt.Sprintf(prompts.DOCTOR, name, restPlayersStr, enums.RoleDoctor)
		} else if i == 3 {
			players[i].Role = enums.RoleDetective
			players[i].SystemPrompt = fmt.Sprintf(prompts.DETECTIVE, name, restPlayersStr, enums.RoleDetective)
		} else {
			players[i].Role = enums.RoleCitizen
			players[i].SystemPrompt = fmt.Sprintf(prompts.CITIZEN, name, restPlayersStr, enums.RoleCitizen)
		}
		players[i].Name = name
	}

	// rand.Shuffle(len(players), func(i, j int) {
	// 	players[i], players[j] = players[j], players[i]
	// })

	llm := llm.GetLLM()

	gameState := state.NewGameState(players, llm)
	status := enums.GameStatusOngoing

	firstDay := true

	for status == enums.GameStatusOngoing {
		if err := gameState.DayPhase(firstDay); err != nil {
			panic(err)
		}
		firstDay = false
		status = gameState.EndgameStatus()
		if status != enums.GameStatusOngoing {
			break
		}
		if err := gameState.NightPhase(); err != nil {
			panic(err)
		}
		gameState.UpdateCycle()
	}

	fmt.Printf("Game ended with status: %d\n", status)
}
