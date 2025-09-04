package secp256k1

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type PrivateKey = ecdsa.PrivateKey
type PublicKey = ecdsa.PublicKey

func GenPubKey(prvkey []byte) []byte {
	x, y := secp256k1.S256().ScalarBaseMult(prvkey)
	return secp256k1.S256().Marshal(x, y)
}

func PrivateKeyFromBytes(prvkey []byte) *PrivateKey {
	var prv PrivateKey
	prv.Curve = secp256k1.S256()
	prv.D = big.NewInt(0).SetBytes(prvkey)

	var pub PublicKey
	pub.Curve = prv.Curve
	pub.X, pub.Y = secp256k1.S256().ScalarBaseMult(prvkey)
	prv.PublicKey = pub
	return &prv
}
