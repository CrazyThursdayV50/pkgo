package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	apiJSON   = jsoniter.ConfigDefault
	Marshal   = apiJSON.Marshal
	Unmarshal = apiJSON.Unmarshal
)

func Init(cfg *jsoniter.Config) {
	apiJSON = cfg.Froze()
	Marshal = apiJSON.Marshal
	Unmarshal = apiJSON.Unmarshal
}

func JSON() jsoniter.API { return apiJSON }
