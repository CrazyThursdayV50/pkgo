package telegram

type Config struct {
	APIKEY        string
	Debug         bool
	Proxy         string
	UpdateFilter  []string
	SkipTlsVerify bool
}
