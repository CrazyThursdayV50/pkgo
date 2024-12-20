package ethereum

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestEthereum(t *testing.T) {
	var prvkey = "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	prvBytes, _ := hex.DecodeString(prvkey)
	var addr = "0x3e9003153d9a39d3f57b126b0c38513d5e289c3e"

	address, err := GenAddressFromPrvKey(prvBytes)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	if address != addr {
		log.Fatalf("addr not equal: %v", address)
	}

	t.Logf("prvkey: %v, address: %v", prvkey, address)
}
