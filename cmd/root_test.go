package cmd

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var canonicalUUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func executeCommand(args ...string) (stdout string, stderr string, err error) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	rootCmd := NewRootCmd(outBuf, errBuf)
	rootCmd.SetArgs(args)

	err = rootCmd.Execute()
	return outBuf.String(), errBuf.String(), err
}

func TestRootCommand(t *testing.T) {
	t.Run("root help output", func(t *testing.T) {
		out, _, err := executeCommand("--help")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "utility-belt") {
			t.Errorf("expected help output to contain 'utility-belt', got %q", out)
		}
		if !strings.Contains(out, "generate") {
			t.Errorf("expected help output to list 'generate' command, got %q", out)
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
