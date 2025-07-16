package main

import (
	"bytes"
	"flag"
	"fmt"
	"mafia/cmd/args"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/game/state"
	"mafia/cmd/llm"
	"math/rand/v2"
	"os"
	"strings"
	"text/template"
	"time"
)

func generateSystemPrompt(tmpl *template.Template, name string, args any) string {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, args); err != nil {
		panic(err)
	}
	return buf.String()
}

func main() {
	out := flag.String("out", "", "Output directory for game logs")
	flag.Parse()

	now := time.Now()

	names := []string{
		"Alice", "Bob", "Charlie", "Diana", "Ethan", "Fiona", "George",
	}

	rand.Shuffle(len(names), func(i, j int) {
		names[i], names[j] = names[j], names[i]
	})

	tmpl := template.Must(template.ParseGlob("templates/*.tmpl"))
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
			players[i].SystemPrompt = generateSystemPrompt(tmpl, "mafia", args.MafiaTemplateArgs{
				CommonTemplateArgs: args.CommonTemplateArgs{
					Name:    name,
					Role:    enums.RoleMafia.String(),
					Players: restPlayersStr,
				},
				Partner: names[(i+1)%2],
			})
		} else if i == 2 {
			players[i].Role = enums.RoleDoctor
			players[i].SystemPrompt = generateSystemPrompt(tmpl, "doctor", args.CommonTemplateArgs{
				Name:    name,
				Role:    enums.RoleDoctor.String(),
				Players: restPlayersStr,
			})
		} else if i == 3 {
			players[i].Role = enums.RoleDetective
			players[i].SystemPrompt = generateSystemPrompt(tmpl, "detective", args.CommonTemplateArgs{
				Name:    name,
				Role:    enums.RoleDetective.String(),
				Players: restPlayersStr,
			})
		} else {
			players[i].Role = enums.RoleCitizen
			players[i].SystemPrompt = generateSystemPrompt(tmpl, "citizen", args.CommonTemplateArgs{
				Name:    name,
				Role:    enums.RoleCitizen.String(),
				Players: restPlayersStr,
			})
		}
		players[i].Name = name
	}

	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	llm := llm.GetOpenaiLLM()

	gameState := state.NewGameState(players, llm)

	firstDay := true

	for gameState.EndgameStatus() == enums.GameStatusOngoing {
		if err := gameState.DayPhase(firstDay); err != nil {
			panic(err)
		}
		firstDay = false
		if gameState.EndgameStatus() != enums.GameStatusOngoing {
			break
		}
		if err := gameState.NightPhase(); err != nil {
			panic(err)
		}
		gameState.UpdateCycle()
	}

	fmt.Printf("Game ended with status: %d\n", gameState.EndgameStatus())

	if out != nil && *out != "" {
		file, err := os.OpenFile(*out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		if err := gameState.Dump(now, players, file); err != nil {
			fmt.Printf("Error dumping game state: %v\n", err)
			return
		}
	}
}
