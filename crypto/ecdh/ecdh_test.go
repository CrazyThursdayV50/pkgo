package ecdh

import (
	"bytes"
	"crypto/ecdh"
	"testing"
)

func TestECDH(t *testing.T) {
	tests := []struct {
		name     string
		curve    ecdh.Curve
		message  []byte
		wantLen  int
		wantFail bool
	}{
		{
			name:    "P256 Test",
			curve:   ecdh.P256(),
			message: []byte("Hello P256"),
			wantLen: 65,
		},
		{
			name:    "P384 Test",
			curve:   ecdh.P384(),
			message: []byte("Hello P384"),
			wantLen: 97,
		},
		{
			name:    "P521 Test",
			curve:   ecdh.P521(),
			message: []byte("Hello P521"),
			wantLen: 133,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建曲线实例
			c := NewCurve(tt.curve)

			// 检查公钥长度
			if c.Len() != tt.wantLen {
				t.Errorf("Len() = %v, want %v", c.Len(), tt.wantLen)
			}

			// 生成密钥对
			privateKey, err := c.RandomEcdh()
			if err != nil {
				t.Fatalf("RandomEcdh() error = %v", err)
			}

			// 加密
			encrypted, err := c.Encrypt(privateKey.PublicKey(), tt.message)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			// 检查加密后的长度
			expectedLen := c.Len() + len(tt.message)
			if len(encrypted) != expectedLen {
				t.Errorf("encrypted length = %v, want %v", len(encrypted), expectedLen)
			}

			// 解密
			decrypted, err := c.Decrypt(privateKey, encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			// 验证解密结果
			if !bytes.Equal(decrypted, tt.message) {
				t.Errorf("Decrypt() = %v, want %v", string(decrypted), string(tt.message))
			}
		})
	}
}

func TestRandomEcdh(t *testing.T) {
	c := NewCurve(nil) // 默认使用 P256
	prv1, err := c.RandomEcdh()
	if err != nil {
		t.Fatalf("RandomEcdh() error = %v", err)
	}

	prv2, err := c.RandomEcdh()
	if err != nil {
		t.Fatalf("RandomEcdh() error = %v", err)
	}

	// 确保生成的两个私钥不相同
	if bytes.Equal(prv1.Bytes(), prv2.Bytes()) {
		t.Error("RandomEcdh() generated identical keys")
	}
}

func TestGenPubKey(t *testing.T) {
	c := NewCurve(nil)
	prv, err := c.RandomEcdh()
	if err != nil {
		t.Fatalf("RandomEcdh() error = %v", err)
	}

	pub, err := c.GenPubKey(prv.Bytes())
	if err != nil {
		t.Fatalf("GenPubKey() error = %v", err)
	}

	// 验证生成的公钥长度
	if len(pub) != c.Len() {
		t.Errorf("GenPubKey() length = %v, want %v", len(pub), c.Len())
	}
}

func TestEncryptDecryptLongMessage(t *testing.T) {
	c := NewCurve(nil)
	prv, err := c.RandomEcdh()
	if err != nil {
		t.Fatalf("RandomEcdh() error = %v", err)
	}

	// 测试长消息
	longMessage := bytes.Repeat([]byte("Long message test "), 100)

	encrypted, err := c.Encrypt(prv.PublicKey(), longMessage)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := c.Decrypt(prv, encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted, longMessage) {
		t.Error("Decrypt() long message failed")
	}
}

func TestInvalidInputs(t *testing.T) {
	c := NewCurve(nil)

	// 测试解密无效数据
	prv, _ := c.RandomEcdh()
	_, err := c.Decrypt(prv, []byte("too short"))
	if err == nil {
		t.Error("Decrypt() should fail with short input")
	}

	// 测试无效的公钥
	invalidPub := make([]byte, c.Len())
	_, err = c.ReadPublic(invalidPub)
	if err == nil {
		t.Error("ReadPublic() should fail with invalid input")
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	c := NewCurve(nil)
	prv, _ := c.RandomEcdh()
	message := []byte("Benchmark test message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, _ := c.Encrypt(prv.PublicKey(), message)
		_, _ = c.Decrypt(prv, encrypted)
	}
}
