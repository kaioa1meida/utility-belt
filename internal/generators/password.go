package generators

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	UppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	NumberChars    = "0123456789"
	SymbolChars    = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
)

var (
	ErrInvalidLength = errors.New("password length must be greater than zero")
	ErrInvalidCount  = errors.New("password count must be greater than zero")
	ErrEmptyCharset  = errors.New("character set cannot be empty: all character classes were excluded")
)

// PasswordOptions defines the configuration for generating passwords.
type PasswordOptions struct {
	Length      int
	Count       int
	NoSymbols   bool
	NoNumbers   bool
	NoUppercase bool
	NoLowercase bool
}

// GeneratePasswords produces one or more cryptographically secure random passwords.
func GeneratePasswords(opts PasswordOptions) ([]string, error) {
	if opts.Length <= 0 {
		return nil, ErrInvalidLength
	}
	if opts.Count <= 0 {
		return nil, ErrInvalidCount
	}

	var charsetBuilder strings.Builder
	if !opts.NoUppercase {
		charsetBuilder.WriteString(UppercaseChars)
	}
	if !opts.NoLowercase {
		charsetBuilder.WriteString(LowercaseChars)
	}
	if !opts.NoNumbers {
		charsetBuilder.WriteString(NumberChars)
	}
	if !opts.NoSymbols {
		charsetBuilder.WriteString(SymbolChars)
	}

	charset := charsetBuilder.String()
	if len(charset) == 0 {
		return nil, ErrEmptyCharset
	}

	charsetLen := big.NewInt(int64(len(charset)))
	passwords := make([]string, opts.Count)

	for i := 0; i < opts.Count; i++ {
		var pwdBuilder strings.Builder
		pwdBuilder.Grow(opts.Length)
		for j := 0; j < opts.Length; j++ {
			idx, err := rand.Int(rand.Reader, charsetLen)
			if err != nil {
				return nil, fmt.Errorf("failed to generate random character: %w", err)
			}
			pwdBuilder.WriteByte(charset[idx.Int64()])
		}
		passwords[i] = pwdBuilder.String()
	}

	return passwords, nil
}
