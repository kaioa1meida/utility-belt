package codecs

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
)

func TestDecodeJWT(t *testing.T) {
	specToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	t.Run("valid spec token", func(t *testing.T) {
		res, err := DecodeJWT(specToken, time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.Header["alg"] != "HS256" {
			t.Errorf("expected alg HS256, got %v", res.Header["alg"])
		}
		if res.Header["typ"] != "JWT" {
			t.Errorf("expected typ JWT, got %v", res.Header["typ"])
		}
		if res.Payload["sub"] != "1234567890" {
			t.Errorf("expected sub 1234567890, got %v", res.Payload["sub"])
		}
		if res.Payload["name"] != "John Doe" {
			t.Errorf("expected name John Doe, got %v", res.Payload["name"])
		}
		if res.IsExpired {
			t.Errorf("expected IsExpired false, got true")
		}
	})

	t.Run("token with future exp is not expired", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour).Unix()
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":"user1","exp":%d}`, future)))
		token := fmt.Sprintf("%s.%s.sig", header, payload)

		res, err := DecodeJWT(token, time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.IsExpired {
			t.Errorf("expected IsExpired false for future exp, got true")
		}
		if res.ExpiresAt == nil || res.ExpiresAt.Unix() != future {
			t.Errorf("expected ExpiresAt %d, got %v", future, res.ExpiresAt)
		}
	})

	t.Run("token with past exp is detected as expired", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour).Unix()
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":"user1","exp":%d}`, past)))
		token := fmt.Sprintf("%s.%s.sig", header, payload)

		res, err := DecodeJWT(token, time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !res.IsExpired {
			t.Errorf("expected IsExpired true for past exp, got false")
		}
	})

	t.Run("invalid token format - wrong number of segments", func(t *testing.T) {
		invalidTokens := []string{
			"",
			"onlyonepart",
			"two.parts",
			"four.parts.are.invalid",
		}

		for _, tok := range invalidTokens {
			_, err := DecodeJWT(tok, time.Now())
			if err == nil {
				t.Errorf("expected error for token %q, got nil", tok)
			}
		}
	})

	t.Run("invalid base64 in header or payload", func(t *testing.T) {
		_, err := DecodeJWT("invalid!base64.eyJzdWIiOiIxMjMifQ.sig", time.Now())
		if err == nil {
			t.Error("expected error for invalid base64 header, got nil")
		}

		_, err = DecodeJWT("eyJhbGciOiJIUzI1NiJ9.invalid!payload.sig", time.Now())
		if err == nil {
			t.Error("expected error for invalid base64 payload, got nil")
		}
	})

	t.Run("invalid json in header or payload", func(t *testing.T) {
		notJSON := base64.RawURLEncoding.EncodeToString([]byte("this is not json"))
		validJSON := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"123"}`))

		_, err := DecodeJWT(fmt.Sprintf("%s.%s.sig", notJSON, validJSON), time.Now())
		if err == nil {
			t.Error("expected error for non-json header, got nil")
		}

		_, err = DecodeJWT(fmt.Sprintf("%s.%s.sig", validJSON, notJSON), time.Now())
		if err == nil {
			t.Error("expected error for non-json payload, got nil")
		}
	})
}
