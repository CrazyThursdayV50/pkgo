package googleauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"math/big"
	"strings"
	"time"

	"github.com/CrazyThursdayV50/pkgo/crypto"
)

var encoder func([]byte) string
var decoder func(string) ([]byte, error)

func init() {
	encoder = base32.StdEncoding.EncodeToString
	decoder = base32.StdEncoding.DecodeString
}

func NewKey() (string, error) {
	key, err := crypto.RandomBytes32()
	if err != nil {
		return "", err
	}

	return encoder(hmacSha1(key, nil)), nil
}

const f = "otpauth://totp/{name}?secret={secret}&issuer={issuer}"

func NewKeyUrl(account string, provider string) (string, error) {
	key, err := NewKey()
	if err != nil {
		return "", err
	}

	var replacer = strings.NewReplacer(
		"{name}", account,
		"{secret}", key,
		"{issuer}", provider,
	)

	return replacer.Replace(f), nil
}

func hmacSha1(key, data []byte) []byte {
	hash := hmac.New(sha1.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

const mask = 0xFF

var (
	shifts = [8]uint8{56, 48, 40, 32, 24, 16, 8, 0}
)

func int64ToBytes(num int64) []byte {
	var data []byte
	for _, shift := range shifts {
		data = append(data, byte((num>>shift)&mask))
	}

	return data
}

func bytesToInt64(data []byte) int64 {
	return int64(data[3])<<0 + int64(data[2])<<8 + int64(data[1])<<16 + int64(data[0])<<24
}

func getIndex(data []byte) byte {
	return data[len(data)-1] & 0x0F
}

func get4Byte(data []byte, from int) []byte {
	return data[from : from+4]
}

func ignoreFirstBit(data []byte) {
	data[0] = data[0] & 0x7F
}

func GenOnetimePassword(key string, ts int64) (int64, error) {
	secret, err := decoder(key)
	if err != nil {
		return 0, err
	}

	hash := hmacSha1(secret, int64ToBytes(ts/30))
	index := getIndex(hash)
	slice := get4Byte(hash, int(index))
	ignoreFirstBit(slice)
	num := big.NewInt(bytesToInt64(slice))
	n := big.NewInt(0).Mod(num, big.NewInt(1e6)).Int64()
	return n, nil
}

func Validate(key string, num int64) bool {
	now := time.Now().Unix()

	if n, err := GenOnetimePassword(key, now); err == nil && n == num {
		return true
	}

	if n, err := GenOnetimePassword(key, now-30); err == nil && n == num {
		return true
	}

	if n, err := GenOnetimePassword(key, now+30); err == nil && n == num {
		return true
	}

	return false
}
