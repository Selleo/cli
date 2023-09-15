package awscmd

import (
	"github.com/Selleo/cli/cryptographic"
	"github.com/Selleo/cli/fmtx"
)

const (
	sesPasswordDate     = "11111111"
	sesPasswordService  = "ses"
	sesPasswordMessage  = "SendRawEmail"
	sesPasswordTerminal = "aws4_request"
	sesPasswordVersion  = 0x04
)

func SESPasswordFromAccessKey(region string, secretKey string) string {
	signature := cryptographic.HMACWithSHA256([]byte("AWS4"+secretKey), []byte(sesPasswordDate))
	signature = cryptographic.HMACWithSHA256(signature, []byte(region))
	signature = cryptographic.HMACWithSHA256(signature, []byte(sesPasswordService))
	signature = cryptographic.HMACWithSHA256(signature, []byte(sesPasswordTerminal))
	signature = cryptographic.HMACWithSHA256(signature, []byte(sesPasswordMessage))
	signatureAndVersion := []byte{sesPasswordVersion}
	signatureAndVersion = append(signatureAndVersion, signature...)
	smtpPassword := fmtx.Base64(signatureAndVersion)
	return smtpPassword
}
