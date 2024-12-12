package tron

import (
	"log"
	"testing"
)

func TestTron(t *testing.T) {
	prvkey, address, err := GenAddress()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	t.Logf("prvkey: %v, address: %v", prvkey, address)
}
