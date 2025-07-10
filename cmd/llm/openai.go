package llm

import (
	"context"
	"fmt"
	"log"
	"mafia/cmd/llm/models"

	"github.com/sashabaranov/go-openai"
)

type OpenaiLLM struct {
	model  string
	client *openai.Client
}

func GetOpenaiLLM() *OpenaiLLM {
	config := ParseConfigArgs()
	client_config := openai.DefaultConfig(config.APIKey)
	client_config.BaseURL = config.BaseURL
	return &OpenaiLLM{
		model:  config.Model,
		client: openai.NewClientWithConfig(client_config),
	}
}

func (o *OpenaiLLM) Generate(ctx context.Context, messages []models.GenerateMessage) (out models.GenerateResponse, err error) {
	chatRequest := openai.ChatCompletionRequest{
		Model: o.model,
	}
	for i := range messages {
		msg := messages[i]
		chatRequest.Messages = append(chatRequest.Messages, openai.ChatCompletionMessage{
			Role:             msg.Role,
			Content:          msg.Content,
			Name:             msg.Name,
			ReasoningContent: msg.ReasoningContent,
		})
	}

	response, err := o.client.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		return out, err
	}

	if len(response.Choices) == 0 {
		return out, fmt.Errorf("no choices returned from OpenAI API")
	}

	if response.Choices[0].Message.ReasoningContent != "" {
		log.Printf("Reasoning: %s\n", response.Choices[0].Message.ReasoningContent)
	}

	out.Content = response.Choices[0].Message.Content
	out.Reasoning = response.Choices[0].Message.ReasoningContent

	return out, nil
}
