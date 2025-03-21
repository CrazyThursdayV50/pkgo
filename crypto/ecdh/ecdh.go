package ecdh

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type Curve struct {
	curve ecdh.Curve
}

func NewCurve(curve ecdh.Curve) *Curve {
	var c Curve
	if curve == nil {
		curve = ecdh.P256()
	}

	c.curve = curve
	return &c
}

func (c *Curve) RandomEcdh() (*ecdh.PrivateKey, error) {
	return c.curve.GenerateKey(rand.Reader)
}

func (c *Curve) GenPubKey(prvkey []byte) ([]byte, error) {
	buf := bytes.NewReader(prvkey)
	private, err := c.curve.GenerateKey(buf)
	if err != nil {
		return nil, err
	}

	return private.PublicKey().Bytes(), nil
}

func (c *Curve) ReadPrivate(key []byte) (*ecdh.PrivateKey, error) {
	return c.curve.NewPrivateKey(key)
}

func (c *Curve) ReadPublic(key []byte) (*ecdh.PublicKey, error) {
	return c.curve.NewPublicKey(key)
}

func (c *Curve) GenSharedPoint(prvA, prvB []byte) ([]byte, error) {
	privateA, err := c.curve.NewPrivateKey(prvA)
	if err != nil {
		return nil, err
	}

	privateB, err := c.curve.NewPrivateKey(prvB)
	if err != nil {
		return nil, err
	}

	return privateA.ECDH(privateB.PublicKey())
}

// Len 返回当前曲线的公钥长度
func (c *Curve) Len() int {
	switch c.curve {
	case ecdh.P256():
		return 65
	case ecdh.P384():
		return 97
	case ecdh.P521():
		return 133
	default:
		return 65 // 默认返回 P256 的长度
	}
}

// Encrypt 加密函数
func (c *Curve) Encrypt(publicKey *ecdh.PublicKey, message []byte) ([]byte, error) {
	// 生成临时密钥对
	ephemeral, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// 计算共享密钥
	sharedSecret, err := ephemeral.ECDH(publicKey)
	if err != nil {
		return nil, err
	}

	// 使用共享密钥派生加密密钥
	h := sha256.New()
	h.Write(sharedSecret)
	key := h.Sum(nil)

	// 加密消息
	ciphertext := make([]byte, len(message))
	for i := 0; i < len(message); i++ {
		ciphertext[i] = message[i] ^ key[i%sha256.Size]
	}

	// 拼接临时公钥和密文
	ephemeralPubKeyBytes := ephemeral.PublicKey().Bytes()
	result := make([]byte, len(ephemeralPubKeyBytes)+len(ciphertext))
	copy(result, ephemeralPubKeyBytes)
	copy(result[len(ephemeralPubKeyBytes):], ciphertext)

	return result, nil
}

// Decrypt 解密函数
func (c *Curve) Decrypt(privateKey *ecdh.PrivateKey, combined []byte) ([]byte, error) {
	// 获取公钥长度
	pubKeyLen := c.Len()
	if len(combined) < pubKeyLen {
		return nil, fmt.Errorf("invalid input length")
	}

	// 分离临时公钥和密文
	ephemeralPubKeyBytes := combined[:pubKeyLen]
	ciphertext := combined[pubKeyLen:]

	// 从字节数组恢复临时公钥
	ephemeralPubKey, err := c.curve.NewPublicKey(ephemeralPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %v", err)
	}

	// 计算共享密钥
	sharedSecret, err := privateKey.ECDH(ephemeralPubKey)
	if err != nil {
		return nil, err
	}

	// 使用共享密钥派生解密密钥
	h := sha256.New()
	h.Write(sharedSecret)
	key := h.Sum(nil)

	// 解密消息
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i++ {
		plaintext[i] = ciphertext[i] ^ key[i%sha256.Size]
	}

	return plaintext, nil
}
