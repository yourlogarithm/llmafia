package llm

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	BaseURL   string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	Model     string `json:"model" yaml:"model"`
	APIKey    string `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	MaxTokens int    `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
}

func ParseConfigArgs() (config Config) {

	base_url := flag.String("base-url", "", "Base URL for the LLM API")
	model := flag.String("model", "", "LLM model to use")
	apiKey := flag.String("api-key", "", "API key for the LLM provider (if required)")
	maxTokens := flag.Int("max-tokens", 8192, "Maximum number of tokens to generate")
	flag.Parse()

	if base_url == nil || *base_url == "" {
		fmt.Println("--base-url must be specified")
		flag.Usage()
		os.Exit(1)
	}

	if model == nil || *model == "" {
		fmt.Println("--model must be specified")
		flag.Usage()
		os.Exit(1)
	}

	config.BaseURL = *base_url
	config.Model = *model
	config.APIKey = *apiKey
	config.MaxTokens = *maxTokens

	if config.APIKey == "" {
		config.APIKey = os.Getenv("API_KEY")
	}

	return config
}
