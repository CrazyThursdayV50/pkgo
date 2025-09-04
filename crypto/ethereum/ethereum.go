package ethereum

import (
	"encoding/hex"

	"github.com/CrazyThursdayV50/pkgo/crypto"
	"github.com/CrazyThursdayV50/pkgo/crypto/ecdsa/secp256k1"
	"golang.org/x/crypto/sha3"
)

func GenAddressFromPubKey(pubkey []byte) (string, error) {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubkey[1:])
	result := hash.Sum(nil)
	return "0x" + hex.EncodeToString(result[12:]), nil
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
