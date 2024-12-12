package sha

import "crypto/sha1"

func Sha1(src []byte) []byte {
	hash := sha1.New()
	hash.Write(src)
	return hash.Sum(nil)
}
