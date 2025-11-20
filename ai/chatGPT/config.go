package chatgpt

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type Config struct {
	BaseURL             string
	Model               string
	Token               string
	MaxTokens           int
	MaxCompletionTokens int
	Temperature         float32
	SystemContent       string
}

func DefaultConfig() *Config {
	return &Config{
		Model:         openai.GPT4Turbo,
		Token:         os.Getenv("OPENAI_API_KEY"),
		MaxTokens:     4000,
		Temperature:   0.5,
		SystemContent: "",
	}
}
