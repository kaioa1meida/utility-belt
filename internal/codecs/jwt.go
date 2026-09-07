package codecs

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidJWTFormat  = errors.New("invalid jwt format: token must have exactly 3 parts separated by dot")
	ErrInvalidJWTHeader  = errors.New("failed to decode jwt header")
	ErrInvalidJWTPayload = errors.New("failed to decode jwt payload")
)

// JWTResult contains the unmarshaled header and payload from a decoded JWT.
type JWTResult struct {
	Header    map[string]any `json:"header"`
	Payload   map[string]any `json:"payload"`
	IsExpired bool           `json:"-"`
	ExpiresAt *time.Time     `json:"-"`
}

func decodeBase64URL(seg string) ([]byte, error) {
	// Standard JWT uses Base64URL without padding (RFC 7515)
	if b, err := base64.RawURLEncoding.DecodeString(seg); err == nil {
		return b, nil
	}
	// Fallback to padded Base64URL
	if b, err := base64.URLEncoding.DecodeString(seg); err == nil {
		return b, nil
	}
	// Fallback to standard base64 if standard padding is used
	return base64.StdEncoding.DecodeString(seg)
}

// DecodeJWT parses and decodes a JWS compact token into Header and Payload.
// It does NOT verify the cryptographic signature.
// If the token contains an "exp" claim and now is after that timestamp, IsExpired is set to true.
func DecodeJWT(tokenStr string, now time.Time) (*JWTResult, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidJWTFormat
	}

	headerBytes, err := decodeBase64URL(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJWTHeader, err)
	}

	payloadBytes, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJWTPayload, err)
	}

	var header map[string]any
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("%w: not valid JSON (%v)", ErrInvalidJWTHeader, err)
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("%w: not valid JSON (%v)", ErrInvalidJWTPayload, err)
	}

	result := &JWTResult{
		Header:  header,
		Payload: payload,
	}

	if expVal, exists := payload["exp"]; exists {
		var expUnix int64
		switch v := expVal.(type) {
		case float64:
			expUnix = int64(v)
		case int64:
			expUnix = v
		case int:
			expUnix = int64(v)
		case json.Number:
			if parsed, err := v.Int64(); err == nil {
				expUnix = parsed
			}
		}

		if expUnix > 0 {
			expTime := time.Unix(expUnix, 0)
			result.ExpiresAt = &expTime
			if !now.IsZero() && now.After(expTime) {
				result.IsExpired = true
			}
		}
	}

	return result, nil
}
