package generators

import (
	"errors"
	"regexp"
	"testing"

	"github.com/google/uuid"
)

var canonicalUUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestGenerateUUIDs(t *testing.T) {
	tests := []struct {
		name       string
		opts       UUIDOptions
		wantErr    error
		validateFn func(t *testing.T, uuids []string)
	}{
		{
			name: "default options (UUID v4, count 1)",
			opts: UUIDOptions{
				V7:    false,
				Count: 1,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, uuids []string) {
				if len(uuids) != 1 {
					t.Fatalf("expected 1 UUID, got %d", len(uuids))
				}
				val := uuids[0]
				if !canonicalUUIDRegex.MatchString(val) {
					t.Errorf("UUID %q does not match canonical format", val)
				}
				parsed, err := uuid.Parse(val)
				if err != nil {
					t.Fatalf("failed to parse generated UUID: %v", err)
				}
				if parsed.Version() != 4 {
					t.Errorf("expected UUID version 4, got %d", parsed.Version())
				}
				if parsed.Variant() != uuid.RFC4122 {
					t.Errorf("expected RFC4122 variant, got %v", parsed.Variant())
				}
			},
		},
		{
			name: "UUID v7 (RFC 9562, count 1)",
			opts: UUIDOptions{
				V7:    true,
				Count: 1,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, uuids []string) {
				if len(uuids) != 1 {
					t.Fatalf("expected 1 UUID, got %d", len(uuids))
				}
				val := uuids[0]
				if !canonicalUUIDRegex.MatchString(val) {
					t.Errorf("UUID %q does not match canonical format", val)
				}
				parsed, err := uuid.Parse(val)
				if err != nil {
					t.Fatalf("failed to parse generated UUID: %v", err)
				}
				if parsed.Version() != 7 {
					t.Errorf("expected UUID version 7, got %d", parsed.Version())
				}
				if parsed.Variant() != uuid.RFC4122 {
					t.Errorf("expected RFC4122 variant, got %v", parsed.Variant())
				}
			},
		},
		{
			name: "batch generation uniqueness (count 100)",
			opts: UUIDOptions{
				V7:    false,
				Count: 100,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, uuids []string) {
				if len(uuids) != 100 {
					t.Fatalf("expected 100 UUIDs, got %d", len(uuids))
				}
				seen := make(map[string]struct{}, len(uuids))
				for _, u := range uuids {
					if !canonicalUUIDRegex.MatchString(u) {
						t.Errorf("UUID %q is not canonical format", u)
					}
					if _, exists := seen[u]; exists {
						t.Fatalf("duplicate UUID generated in batch: %s", u)
					}
					seen[u] = struct{}{}
				}
			},
		},
		{
			name: "batch generation v7 uniqueness (count 100)",
			opts: UUIDOptions{
				V7:    true,
				Count: 100,
			},
			wantErr: nil,
			validateFn: func(t *testing.T, uuids []string) {
				if len(uuids) != 100 {
					t.Fatalf("expected 100 UUIDs, got %d", len(uuids))
				}
				seen := make(map[string]struct{}, len(uuids))
				for _, u := range uuids {
					if !canonicalUUIDRegex.MatchString(u) {
						t.Errorf("UUID %q is not canonical format", u)
					}
					if _, exists := seen[u]; exists {
						t.Fatalf("duplicate UUID generated in batch: %s", u)
					}
					seen[u] = struct{}{}
				}
			},
		},
		{
			name: "invalid count zero",
			opts: UUIDOptions{
				Count: 0,
			},
			wantErr: ErrInvalidUUIDCount,
		},
		{
			name: "invalid count negative",
			opts: UUIDOptions{
				Count: -10,
			},
			wantErr: ErrInvalidUUIDCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateUUIDs(tt.opts)
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
