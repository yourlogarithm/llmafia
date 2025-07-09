package llm

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Provider  string `json:"provider" yaml:"provider"`
	Model     string `json:"model" yaml:"model"`
	APIKey    string `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	MaxTokens int    `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
}

func ParseConfigArgs() (config Config) {
	provider := flag.String("provider", "", "LLM provider (e.g., ollama, openai)")
	model := flag.String("model", "", "LLM model to use")
	apiKey := flag.String("api_key", "", "API key for the LLM provider (if required)")
	maxTokens := flag.Int("max_tokens", 8192, "Maximum number of tokens to generate")
	flag.Parse()

	if provider == nil || *provider == "" {
		fmt.Println("--provider must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if model == nil || *model == "" {
		fmt.Println("--model must be specified")
		flag.Usage()
		os.Exit(1)
	}

	config.Model = *model
	config.Provider = *provider
	config.APIKey = *apiKey
	config.MaxTokens = *maxTokens

	if config.APIKey == "" {
		config.APIKey = os.Getenv("API_KEY")
	}

	return config
}
