package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func EncryptGCM(key []byte, data []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func DecryptGCM(key []byte, encrypted []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := encrypted[:gcm.NonceSize()], encrypted[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
