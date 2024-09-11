package chatgpt

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type Config struct {
	Model       string
	Token       string
	MaxTokens   int
	Temperature float32
	SystemFile  string
}

func DefaultConfig() *Config {
	return &Config{
		Model:       openai.GPT4Turbo,
		Token:       os.Getenv("OPENAI_API_KEY"),
		MaxTokens:   4000,
		Temperature: 0.5,
		SystemFile:  "",
	}
}
