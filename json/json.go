package json

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

var (
	apiJSON   = jsoniter.ConfigDefault
	Marshal   = apiJSON.Marshal
	Unmarshal = apiJSON.Unmarshal
)

type RawMessage = json.RawMessage

func Init(cfg *jsoniter.Config) {
	apiJSON = cfg.Froze()
	Marshal = apiJSON.Marshal
	Unmarshal = apiJSON.Unmarshal
}

func JSON() jsoniter.API { return apiJSON }
