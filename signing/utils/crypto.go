package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMacSha256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func HexSha256(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
