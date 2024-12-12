package crypto

import (
	"crypto/rand"
)

// 生成随机数
// len: 控制随机数的长度
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RandomBytes32() ([]byte, error) {
	return RandomBytes(32)
}
