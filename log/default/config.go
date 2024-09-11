package defaultlogger

type Config struct {
	Mode              string
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
	CallerSkip        int
}

func DefaultConfig() *Config {
	return &Config{
		Mode:              "Development",
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "console",
		Level:             "debug",
		CallerSkip:        1,
	}
}
