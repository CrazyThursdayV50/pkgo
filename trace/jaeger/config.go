package jaeger

type Config struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

func DefaultConfig() *Config {
	return &Config{
		Host:        "localhost:14268",
		ServiceName: "DefaultService",
		LogSpans:    true,
	}

}
