package fmtx

import (
	"encoding/base32"
	"encoding/base64"
)

func Base64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base32(b []byte) string {
	return base32.StdEncoding.EncodeToString(b)
}
