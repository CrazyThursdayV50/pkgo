package tron

import (
	"github.com/CrazyThursdayV50/pkgo/crypto"
	"github.com/CrazyThursdayV50/pkgo/crypto/ecdsa/secp256k1"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"
)

const version = 0x41

func GenAddressFromPubKey(pubkey []byte) (string, error) {
	hash := sha3.New256()
	hash.Write(pubkey)
	result := hash.Sum(nil)
	l := len(result)
	input := result[l-1-20 : l-1]
	return base58.CheckEncode(input, version), nil
}

func GenAddressFromPrvKey(prvkey []byte) (string, error) {
	pubkey := secp256k1.GenPubKey(prvkey)
	addr, err := GenAddressFromPubKey(pubkey)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func GenAddress() ([]byte, string, error) {
	prvkey, err := crypto.RandomBytes32()
	if err != nil {
		return nil, "", err
	}

	addr, err := GenAddressFromPrvKey(prvkey)
	if err != nil {
		return nil, "", err
	}
	return prvkey, addr, nil
}
