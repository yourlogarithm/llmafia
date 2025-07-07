package state

import (
	"context"
	"fmt"
	"mafia/cmd/enums"
	"mafia/cmd/game"
	"mafia/cmd/utils"
	"strings"

	"github.com/teilomillet/gollm"
)

func (gs *GameState) dayDiscussion(player *game.Player) error {
	narratorText := fmt.Sprintf("Now it is %s's turn to discuss.", player.Name)
	gs.Conversation.AddMessage(
		game.NARRATOR,
		narratorText,
	)
	prompt := gs.basePrompt(player)

	response, err := gs.llm.Generate(context.Background(), &prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	gs.Conversation.AddMessage(
		player,
		strings.Trim(response, " \n"),
	)

	if err := gs.SolicitVote(player); err != nil {
		return fmt.Errorf("failed to solicit vote: %w", err)
	}

	return nil
}

func (gs *GameState) SolicitVote(player *game.Player) error {
	prompt := gs.basePrompt(player)

	vote_prompt := "You may optionally accuse someone of being a mafia member. In case you do so, at the end of the day, all players will vote to eliminate the accused player. If the majority of players vote to eliminate the accused player, he or she will be eliminated. The response format must be JSON:\n{\"accuse\": \"<player_name>\", \"reason\": \"explain why you <player_name>\"}\n\nIf you do not want to accuse anyone, just leave `accuse` and `reason` empty: {\"accuse\": \"\", \"reason\": \"\"}."

	prompt.Messages = append(prompt.Messages, gollm.PromptMessage{
		Role:    "user",
		Name:    string(enums.RoleNarrator),
		Content: vote_prompt,
	})

	var accuse struct {
		Accuse string `json:"accuse"`
		Reason string `json:"reason"`
	}

	response, err := gs.llm.Generate(context.Background(), &prompt)
	if err != nil {
		return fmt.Errorf("failed to generate response for accusation: %w", err)
	}

	if err := utils.ParseJSONResponsePermissive(response, &accuse); err != nil {
		return fmt.Errorf("failed to parse accusation response: %w", err)
	}

	if accuse.Accuse != "" {
		for _, p := range gs.players {
			if p.Name == accuse.Accuse {
				gs.accusedPlayers[p.Name] = player.Name
				gs.Conversation.AddMessage(
					player,
					fmt.Sprintf("I accuse %s of being a Mafia member, therefore a vote for elimination will be proposed at the end of the day. %s", p.Name, accuse.Reason),
				)
			}
		}
	}

	return nil
}
