package ecdsa

import "github.com/ethereum/go-ethereum/crypto/secp256k1"

func GenPubKey(prvkey []byte) []byte {
	x, y := secp256k1.S256().ScalarBaseMult(prvkey)
	return secp256k1.S256().Marshal(x, y)
}
