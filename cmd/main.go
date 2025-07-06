package main

import (
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/game/state"
	"mafia/cmd/llm"
	"mafia/cmd/prompts"
)

func main() {
	names := []string{
		"Alice", "Bob", "Charlie", "Diana", "Ethan", "Fiona", "George", "Hannah",
	}
	// names = names[:2]
	// rand.Shuffle(len(names), func(i, j int) {
	// 	names[i], names[j] = names[j], names[i]
	// })

	players := make([]game.Player, len(names))
	for i, name := range names {
		if i < 2 {
			players[i].Role = enums.RoleMafia
			players[i].SystemPrompt = fmt.Sprintf(prompts.MAFIA, name, enums.RoleMafia, names[(i+1)%2])
		} else if i == 2 {
			players[i].Role = enums.RoleDoctor
			players[i].SystemPrompt = fmt.Sprintf(prompts.DOCTOR, name, enums.RoleDoctor)
		} else if i == 3 {
			players[i].Role = enums.RoleDetective
			players[i].SystemPrompt = fmt.Sprintf(prompts.DETECTIVE, name, enums.RoleDetective)
		} else {
			players[i].Role = enums.RoleCitizen
			players[i].SystemPrompt = fmt.Sprintf(prompts.CITIZEN, name, enums.RoleCitizen)
		}
		players[i].Name = name
	}

	// rand.Shuffle(len(players), func(i, j int) {
	// 	players[i], players[j] = players[j], players[i]
	// })

	llm := llm.GetLLM()

	gameState := state.NewGameState(players, llm)

	err := gameState.DayPhase()
	if err != nil {
		panic(err)
	}
}
