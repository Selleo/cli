package random

import (
	"crypto/rand"

	"github.com/Selleo/cli/fmtx"
)

func Bytes(size int, format string) (string, error) {
	b := make([]byte, size)
	_, _ = rand.Read(b)
	return fmtx.OutputFormat(b, format)
}
