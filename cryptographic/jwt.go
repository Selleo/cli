package cryptographic

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"strings"
)

type JWT struct {
	Header    map[string]any
	Payload   map[string]any
	Signature []byte
	// Verified is true only when a secret was given and the signature matched.
	Verified bool
}

var jwtB64 = base64.RawURLEncoding

func jwtHash(alg string) (func() hash.Hash, error) {
	switch strings.ToUpper(alg) {
	case "HS256":
		return sha256.New, nil
	case "HS384":
		return sha512.New384, nil
	case "HS512":
		return sha512.New, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s (supported: HS256, HS384, HS512)", alg)
	}
}

func jwtSign(alg string, secret []byte, signingInput string) ([]byte, error) {
	h, err := jwtHash(alg)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(h, secret)
	mac.Write([]byte(signingInput))
	return mac.Sum(nil), nil
}

// JWTEncode signs the JSON object payload with the given HMAC algorithm.
func JWTEncode(alg string, secret []byte, payload map[string]any) (string, error) {
	alg = strings.ToUpper(alg)
	if _, err := jwtHash(alg); err != nil {
		return "", err
	}

	header, err := json.Marshal(map[string]string{"alg": alg, "typ": "JWT"})
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	signingInput := jwtB64.EncodeToString(header) + "." + jwtB64.EncodeToString(body)
	sig, err := jwtSign(alg, secret, signingInput)
	if err != nil {
		return "", err
	}
	return signingInput + "." + jwtB64.EncodeToString(sig), nil
}

// JWTDecode parses the token. When secret is non-nil the signature is verified
// and an error is returned if it does not match.
func JWTDecode(token string, secret []byte) (*JWT, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token: expected 3 dot-separated segments")
	}

	jwt := &JWT{}
	if err := jwtDecodeSegment(parts[0], &jwt.Header); err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}
	if err := jwtDecodeSegment(parts[1], &jwt.Payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}
	sig, err := jwtB64.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}
	jwt.Signature = sig

	if secret != nil {
		alg, _ := jwt.Header["alg"].(string)
		expected, err := jwtSign(alg, secret, parts[0]+"."+parts[1])
		if err != nil {
			return nil, err
		}
		if !hmac.Equal(sig, expected) {
			return nil, errors.New("invalid signature")
		}
		jwt.Verified = true
	}

	return jwt, nil
}

func jwtDecodeSegment(seg string, v any) error {
	raw, err := jwtB64.DecodeString(seg)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return dec.Decode(v)
}
