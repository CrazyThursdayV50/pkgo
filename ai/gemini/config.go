package gemini

import "google.golang.org/genai"

type Config struct {
	Token      string
	Model      string
	SystemFile string

	thinkingCfg *genai.ThinkingConfig
}

func (c *Config) SetThinkingConfig(cfg *genai.ThinkingConfig) {
	c.thinkingCfg = cfg
}
