package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

const (
	rsaBits2048 = 2 << 10
	rsaBits4096 = 4 << 10
)

func NewRsaPrivateKey(bit int) (pk *rsa.PrivateKey) {
	var err error
	switch bit {
	case rsaBits2048, rsaBits4096:
		pk, err = rsa.GenerateKey(rand.Reader, bit)
	default:
		pk, err = rsa.GenerateKey(rand.Reader, rsaBits2048)
	}
	if err != nil {
		panic(err)
	}
	return
}

func RsaPrivateKeyToHex(prv *rsa.PrivateKey) string {
	bytes, err := x509.MarshalPKCS8PrivateKey(prv)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
