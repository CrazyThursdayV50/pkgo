package zap

type Config struct {
	Development       bool
	Console           bool
	DisableCaller     bool
	DisableStacktrace bool
	Level             string
	CallerSkip        int
}

func DefaultConfig() *Config {
	return &Config{
		Development:       true,
		Console:           true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Level:             "debug",
		CallerSkip:        0,
	}
}
