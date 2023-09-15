package cryptographic

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HMACWithSHA256(key, message []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return expectedMAC
}
