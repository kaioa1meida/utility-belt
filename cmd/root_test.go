package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

var canonicalUUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func executeCommandWithInput(input string, args ...string) (stdout string, stderr string, err error) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	rootCmd := NewRootCmd(outBuf, errBuf)
	if input != "" {
		rootCmd.SetIn(bytes.NewBufferString(input))
	}
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func executeCommand(args ...string) (stdout string, stderr string, err error) {
	return executeCommandWithInput("", args...)
}

func TestRootCommand(t *testing.T) {
	t.Run("root help output lists categories", func(t *testing.T) {
		out, _, err := executeCommand("--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, cat := range []string{"utility-belt", "generate", "encode", "decode"} {
			if !strings.Contains(out, cat) {
				t.Errorf("expected help output to contain %q, got %q", cat, out)
			}
		}
	})

	t.Run("generate help output", func(t *testing.T) {
		out, _, err := executeCommand("generate", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "password") || !strings.Contains(out, "uuid") {
			t.Errorf("expected generate help to list 'password' and 'uuid', got %q", out)
		}
	})

	t.Run("encode help output", func(t *testing.T) {
		out, _, err := executeCommand("encode", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "base64") {
			t.Errorf("expected encode help to list 'base64', got %q", out)
		}
	})

	t.Run("decode help output", func(t *testing.T) {
		out, _, err := executeCommand("decode", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "base64") || !strings.Contains(out, "jwt") {
			t.Errorf("expected decode help to list 'base64' and 'jwt', got %q", out)
		}
	})

	t.Run("unknown command returns error", func(t *testing.T) {
		_, _, err := executeCommand("unknown-command")
		if err == nil {
			t.Fatal("expected error for unknown command, got nil")
		}
	})
}

func TestGeneratePasswordCommand(t *testing.T) {
	t.Run("default password generation", func(t *testing.T) {
		out, errOut, err := executeCommand("generate", "password")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected exactly 1 line, got %d", len(lines))
		}
		if len(lines[0]) != 16 {
			t.Errorf("expected password length 16, got %d (%q)", len(lines[0]), lines[0])
		}
	})

	t.Run("custom length and count", func(t *testing.T) {
		out, errOut, err := executeCommand("generate", "password", "--length", "24", "--count", "4")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 4 {
			t.Fatalf("expected 4 lines, got %d", len(lines))
		}
		for i, line := range lines {
			if len(line) != 24 {
				t.Errorf("line[%d]: expected length 24, got %d (%q)", i, len(line), line)
			}
		}
	})

	t.Run("no-symbols flag", func(t *testing.T) {
		out, _, err := executeCommand("generate", "password", "--length", "50", "--no-symbols")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		pwd := strings.TrimSpace(out)
		for _, ch := range pwd {
			if strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:,.<>?/~", ch) {
				t.Errorf("found symbol %c in password %q", ch, pwd)
			}
		}
	})

	t.Run("empty alphabet returns error", func(t *testing.T) {
		out, _, err := executeCommand(
			"generate", "password",
			"--no-symbols", "--no-numbers", "--no-uppercase", "--no-lowercase",
		)
		if err == nil {
			t.Fatal("expected error when all character sets are excluded, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("invalid length returns error", func(t *testing.T) {
		out, _, err := executeCommand("generate", "password", "--length", "0")
		if err == nil {
			t.Fatal("expected error for --length 0, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("invalid count returns error", func(t *testing.T) {
		out, _, err := executeCommand("generate", "password", "--count", "0")
		if err == nil {
			t.Fatal("expected error for --count 0, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("password help contains examples and flags", func(t *testing.T) {
		out, _, err := executeCommand("generate", "password", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, expected := range []string{"--length", "--count", "--no-symbols", "--no-numbers", "--no-uppercase", "--no-lowercase"} {
			if !strings.Contains(out, expected) {
				t.Errorf("expected password help to contain %q", expected)
			}
		}
	})
}

func TestGenerateUUIDCommand(t *testing.T) {
	t.Run("default UUID v4 generation", func(t *testing.T) {
		out, errOut, err := executeCommand("generate", "uuid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected exactly 1 line, got %d", len(lines))
		}
		val := lines[0]
		if !canonicalUUIDRegex.MatchString(val) {
			t.Fatalf("UUID %q is not canonical format", val)
		}
		parsed, parseErr := uuid.Parse(val)
		if parseErr != nil {
			t.Fatalf("failed to parse UUID: %v", parseErr)
		}
		if parsed.Version() != 4 {
			t.Errorf("expected UUID version 4, got %d", parsed.Version())
		}
	})

	t.Run("UUID v7 generation", func(t *testing.T) {
		out, errOut, err := executeCommand("generate", "uuid", "--v7")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected exactly 1 line, got %d", len(lines))
		}
		val := lines[0]
		if !canonicalUUIDRegex.MatchString(val) {
			t.Fatalf("UUID %q is not canonical format", val)
		}
		parsed, parseErr := uuid.Parse(val)
		if parseErr != nil {
			t.Fatalf("failed to parse UUID: %v", parseErr)
		}
		if parsed.Version() != 7 {
			t.Errorf("expected UUID version 7, got %d", parsed.Version())
		}
	})

	t.Run("batch UUID count", func(t *testing.T) {
		out, errOut, err := executeCommand("generate", "uuid", "--count", "5")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 5 {
			t.Fatalf("expected 5 lines, got %d", len(lines))
		}

		seen := make(map[string]struct{})
		for _, u := range lines {
			if !canonicalUUIDRegex.MatchString(u) {
				t.Errorf("UUID %q is not in canonical format", u)
			}
			if _, exists := seen[u]; exists {
				t.Errorf("duplicate UUID %q found in batch", u)
			}
			seen[u] = struct{}{}
		}
	})

	t.Run("invalid count returns error", func(t *testing.T) {
		out, _, err := executeCommand("generate", "uuid", "--count", "0")
		if err == nil {
			t.Fatal("expected error for --count 0, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("uuid help contains flags and description", func(t *testing.T) {
		out, _, err := executeCommand("generate", "uuid", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, expected := range []string{"--v7", "--count"} {
			if !strings.Contains(out, expected) {
				t.Errorf("expected uuid help to contain %q", expected)
			}
		}
	})
}

func TestEncodeBase64Command(t *testing.T) {
	t.Run("encode via flag --string", func(t *testing.T) {
		out, errOut, err := executeCommand("encode", "base64", "--string", "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}
		if strings.TrimSpace(out) != "aGVsbG8=" {
			t.Errorf("expected aGVsbG8=, got %q", out)
		}
	})

	t.Run("encode via stdin pipe", func(t *testing.T) {
		out, errOut, err := executeCommandWithInput("hello", "encode", "base64")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}
		if strings.TrimSpace(out) != "aGVsbG8=" {
			t.Errorf("expected aGVsbG8=, got %q", out)
		}
	})

	t.Run("encode without input returns error", func(t *testing.T) {
		out, _, err := executeCommand("encode", "base64")
		if err == nil {
			t.Fatal("expected error when no input provided, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("encode help contains flag and description", func(t *testing.T) {
		out, _, err := executeCommand("encode", "base64", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "--string") {
			t.Errorf("expected help to describe --string flag, got %q", out)
		}
	})
}

func TestDecodeBase64Command(t *testing.T) {
	t.Run("decode via flag --string", func(t *testing.T) {
		out, errOut, err := executeCommand("decode", "base64", "--string", "aGVsbG8=")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}
		if strings.TrimSpace(out) != "hello" {
			t.Errorf("expected hello, got %q", out)
		}
	})

	t.Run("decode via stdin pipe", func(t *testing.T) {
		out, errOut, err := executeCommandWithInput("aGVsbG8=", "decode", "base64")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}
		if strings.TrimSpace(out) != "hello" {
			t.Errorf("expected hello, got %q", out)
		}
	})

	t.Run("decode invalid base64 returns error", func(t *testing.T) {
		out, _, err := executeCommandWithInput("invalid!!!", "decode", "base64")
		if err == nil {
			t.Fatal("expected error for invalid base64, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("decode without input returns error", func(t *testing.T) {
		out, _, err := executeCommand("decode", "base64")
		if err == nil {
			t.Fatal("expected error when no input provided, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("decode help contains flag and description", func(t *testing.T) {
		out, _, err := executeCommand("decode", "base64", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "--string") {
			t.Errorf("expected help to describe --string flag, got %q", out)
		}
	})
}

func TestDecodeJWTCommand(t *testing.T) {
	specToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	t.Run("decode valid JWT via flag --string", func(t *testing.T) {
		out, errOut, err := executeCommand("decode", "jwt", "--string", specToken)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		var parsed struct {
			Header  map[string]any `json:"header"`
			Payload map[string]any `json:"payload"`
		}
		if err := json.Unmarshal([]byte(out), &parsed); err != nil {
			t.Fatalf("failed to unmarshal JSON output %q: %v", out, err)
		}
		if parsed.Header["alg"] != "HS256" {
			t.Errorf("expected alg HS256, got %v", parsed.Header["alg"])
		}
		if parsed.Payload["sub"] != "1234567890" {
			t.Errorf("expected sub 1234567890, got %v", parsed.Payload["sub"])
		}
	})

	t.Run("decode valid JWT via stdin pipe", func(t *testing.T) {
		out, errOut, err := executeCommandWithInput(specToken, "decode", "jwt")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if errOut != "" {
			t.Errorf("expected empty stderr, got %q", errOut)
		}

		if !strings.Contains(out, `"sub": "1234567890"`) {
			t.Errorf("expected output to contain payload sub, got %q", out)
		}
	})

	t.Run("decode expired JWT prints warning to stderr and json to stdout", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour).Unix()
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":"user1","exp":%d}`, past)))
		token := fmt.Sprintf("%s.%s.sig", header, payload)

		out, errOut, err := executeCommand("decode", "jwt", "--string", token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(errOut, "warning: token is expired") {
			t.Errorf("expected warning in stderr, got %q", errOut)
		}
		if !strings.Contains(out, `"sub": "user1"`) {
			t.Errorf("expected json payload in stdout, got %q", out)
		}
	})

	t.Run("decode invalid JWT format returns error", func(t *testing.T) {
		out, _, err := executeCommand("decode", "jwt", "--string", "invalid.jwt")
		if err == nil {
			t.Fatal("expected error for 2-part jwt, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("decode jwt without input returns error", func(t *testing.T) {
		out, _, err := executeCommand("decode", "jwt")
		if err == nil {
			t.Fatal("expected error when no input provided, got nil")
		}
		if out != "" {
			t.Errorf("expected empty stdout on error, got %q", out)
		}
	})

	t.Run("jwt help contains flag and description", func(t *testing.T) {
		out, _, err := executeCommand("decode", "jwt", "--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "--string") {
			t.Errorf("expected help to describe --string flag, got %q", out)
		}
	})
}
