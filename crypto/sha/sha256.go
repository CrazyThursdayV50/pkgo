package sha

import (
	"crypto/sha256"
)

func Sha256(src []byte) []byte {
	hash := sha256.New()
	hash.Write(src)
	return hash.Sum(nil)
}

func Sha256Double(src []byte) []byte {
	return Sha256(Sha256(src))
}
