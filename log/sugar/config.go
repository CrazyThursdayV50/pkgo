package sugar

import "github.com/CrazyThursdayV50/pkgo/log/zap"

type Config = zap.Config

func DefaultConfig() *Config {
	cfg := zap.DefaultConfig()
	cfg.CallerSkip = 1
	return cfg
}
