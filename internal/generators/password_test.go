package generators

import (
	"errors"
	"strings"
	"testing"
)

func TestGeneratePasswords(t *testing.T) {
	tests := []struct {
		name       string
		opts       PasswordOptions
		wantErr    error
		validateFn func(t *testing.T, pwds []string)
	}{
		{
			name: "default options (16 chars, count 1, all classes enabled)",
			opts: PasswordOptions{
				Length: 16,
				Count:  1,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				if len(pwds) != 1 {
					t.Fatalf("expected 1 password, got %d", len(pwds))
				}
				if len(pwds[0]) != 16 {
					t.Errorf("expected length 16, got %d", len(pwds[0]))
				}
			},
		},
		{
			name: "custom length 32 and count 5",
			opts: PasswordOptions{
				Length: 32,
				Count:  5,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				if len(pwds) != 5 {
					t.Fatalf("expected 5 passwords, got %d", len(pwds))
				}
				for i, p := range pwds {
					if len(p) != 32 {
						t.Errorf("password[%d]: expected length 32, got %d", i, len(p))
					}
				}
			},
		},
		{
			name: "no symbols",
			opts: PasswordOptions{
				Length:    50,
				Count:     3,
				NoSymbols: true,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				for _, p := range pwds {
					if strings.ContainsAny(p, SymbolChars) {
						t.Errorf("expected no symbols in password %q", p)
					}
				}
			},
		},
		{
			name: "no numbers",
			opts: PasswordOptions{
				Length:    50,
				Count:     3,
				NoNumbers: true,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				for _, p := range pwds {
					if strings.ContainsAny(p, NumberChars) {
						t.Errorf("expected no numbers in password %q", p)
					}
				}
			},
		},
		{
			name: "no uppercase",
			opts: PasswordOptions{
				Length:      50,
				Count:       3,
				NoUppercase: true,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				for _, p := range pwds {
					if strings.ContainsAny(p, UppercaseChars) {
						t.Errorf("expected no uppercase letters in password %q", p)
					}
				}
			},
		},
		{
			name: "no lowercase",
			opts: PasswordOptions{
				Length:      50,
				Count:       3,
				NoLowercase: true,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				for _, p := range pwds {
					if strings.ContainsAny(p, LowercaseChars) {
						t.Errorf("expected no lowercase letters in password %q", p)
					}
				}
			},
		},
		{
			name: "numbers only",
			opts: PasswordOptions{
				Length:      20,
				Count:       2,
				NoSymbols:   true,
				NoUppercase: true,
				NoLowercase: true,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, pwds []string) {
				for _, p := range pwds {
					for _, ch := range p {
						if !strings.ContainsRune(NumberChars, ch) {
							t.Errorf("expected only digits, got char %c in %q", ch, p)
						}
					}
				}
			},
		},
		{
			name: "empty charset error when all classes are excluded",
			opts: PasswordOptions{
				Length:      16,
				Count:       1,
				NoSymbols:   true,
				NoNumbers:   true,
				NoUppercase: true,
				NoLowercase: true,
			},
			wantErr: ErrEmptyCharset,
		},
		{
			name: "invalid length zero",
			opts: PasswordOptions{
				Length: 0,
				Count:  1,
			},
			wantErr: ErrInvalidLength,
		},
		{
			name: "invalid length negative",
			opts: PasswordOptions{
				Length: -5,
				Count:  1,
			},
			wantErr: ErrInvalidLength,
		},
		{
			name: "invalid count zero",
			opts: PasswordOptions{
				Length: 16,
				Count:  0,
			},
			wantErr: ErrInvalidCount,
		},
		{
			name: "invalid count negative",
			opts: PasswordOptions{
				Length: 16,
				Count:  -2,
			},
			wantErr: ErrInvalidCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePasswords(tt.opts)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateFn != nil {
				tt.validateFn(t, got)
			}
		})
	}
}
