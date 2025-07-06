package llm

import (
	"time"

	"github.com/teilomillet/gollm"
)

func GetLLM() gollm.LLM {
	llm, err := gollm.NewLLM(
		gollm.SetProvider("ollama"),
		gollm.SetModel("gemma3n:e4b"),
		// gollm.SetAPIKey(apiKey),
		// gollm.SetMaxTokens(200),
		gollm.SetMaxRetries(3),
		gollm.SetRetryDelay(time.Second*2),
		gollm.SetLogLevel(gollm.LogLevelInfo),
	)
	if err != nil {
		panic(err)
	}
	return llm
}
