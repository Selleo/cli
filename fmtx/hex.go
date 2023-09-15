package fmtx

import "encoding/hex"

func Hex(b []byte) string {
	return hex.EncodeToString(b)
}
