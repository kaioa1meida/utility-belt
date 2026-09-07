package codecs

import (
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrEmptyBase64Input = errors.New("cannot decode empty base64 input")
)

// EncodeBase64 encodes raw byte input into standard RFC 4648 Base64 string.
func EncodeBase64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// DecodeBase64 decodes a standard Base64 string back into raw bytes.
// Leading and trailing whitespace is stripped before decoding.
func DecodeBase64(input string) ([]byte, error) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 {
		return nil, ErrEmptyBase64Input
	}
	return base64.StdEncoding.DecodeString(trimmed)
}
