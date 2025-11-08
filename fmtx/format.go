package fmtx

import (
	"fmt"
	"io"

	"github.com/wzshiming/ctc"
)

func OutputFormat(b []byte, format string) (string, error) {
	switch format {
	case "hex":
		return Hex(b), nil
	case "base64":
		return Base64(b), nil
	case "base32":
		return Base32(b), nil
	case "raw":
		return RawBytes(b), nil
	}

	return "", fmt.Errorf("unknown format: %s", format)
}

func FGreenln(w io.Writer, s string) {
	fmt.Fprintf(w, "%s%s%s\n", ctc.ForegroundGreen, s, ctc.Reset)
}

func FYellowln(w io.Writer, s string) {
	fmt.Fprintf(w, "%s%s%s\n", ctc.ForegroundYellow, s, ctc.Reset)
}
