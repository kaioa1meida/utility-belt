package codecs

import (
	"bytes"
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple text hello",
			input:    []byte("hello"),
			expected: "aGVsbG8=",
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: "",
		},
		{
			name:     "symbols and numbers",
			input:    []byte("12345!@#$%^&*()_+"),
			expected: "MTIzNDUhQCMkJV4mKigpXys=",
		},
		{
			name:     "binary bytes",
			input:    []byte{0x00, 0xFF, 0xFE, 0x01},
			expected: "AP/+AQ==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeBase64(tt.input)
			if got != tt.expected {
				t.Errorf("EncodeBase64() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []byte
		expectErr bool
	}{
		{
			name:      "valid base64 hello",
			input:     "aGVsbG8=",
			expected:  []byte("hello"),
			expectErr: false,
		},
		{
			name:      "valid base64 with leading and trailing newlines",
			input:     "\naGVsbG8=\n",
			expected:  []byte("hello"),
			expectErr: false,
		},
		{
			name:      "valid base64 with spaces",
			input:     "  aGVsbG8=  ",
			expected:  []byte("hello"),
			expectErr: false,
		},
		{
			name:      "empty input",
			input:     "",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "whitespace only",
			input:     "   \n\t  ",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "invalid base64 characters",
			input:     "invalid!!!",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "invalid padding",
			input:     "aGVsbG",
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase64(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("DecodeBase64(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("DecodeBase64(%q) unexpected error: %v", tt.input, err)
			}
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("DecodeBase64(%q) = %q, want %q", tt.input, string(got), string(tt.expected))
			}
		})
	}
}
