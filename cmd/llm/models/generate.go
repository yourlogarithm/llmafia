package models

type GenerateMessage struct {
	Role             string
	Content          string
	ReasoningContent string
	Name             string
}

type GenerateResponse struct {
	Content   string
	Reasoning string
}
