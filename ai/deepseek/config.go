package deepseek

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type Config struct {
	Model       string
	Token       string
	Temperature float32
	SystemFile  string
}

func DefaultConfig() *Config {
	return &Config{
		Model:       openai.GPT4Turbo,
		Token:       os.Getenv("OPENAI_API_KEY"),
		Temperature: 0.5,
		SystemFile:  "",
	}
}
