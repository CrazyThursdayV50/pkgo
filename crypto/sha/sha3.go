package sha

import "golang.org/x/crypto/sha3"

func Sha3_256(src []byte) []byte {
	hash := sha3.New256()
	hash.Write(src)
	return hash.Sum(nil)
}
