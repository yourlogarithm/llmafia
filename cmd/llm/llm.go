package llm

import (
	"context"
	"mafia/cmd/llm/models"
)

type LLM interface {
	Generate(ctx context.Context, messages []models.GenerateMessage) (models.GenerateResponse, error)
}
