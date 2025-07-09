package llm

import (
	"time"

	"github.com/teilomillet/gollm"
)

func GetLLM() gollm.LLM {
	config := ParseConfigArgs()

	llm, err := gollm.NewLLM(
		gollm.SetProvider(config.Provider),
		gollm.SetModel(config.Model),
		gollm.SetAPIKey(config.APIKey),
		gollm.SetMaxTokens(config.MaxTokens),
		gollm.SetMaxRetries(3),
		gollm.SetRetryDelay(time.Second*2),
		gollm.SetLogLevel(gollm.LogLevelDebug),
	)
	if err != nil {
		panic(err)
	}
	return llm
}
