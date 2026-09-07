package generators

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidUUIDCount = errors.New("uuid count must be greater than zero")
	ErrUUIDCollision    = errors.New("uuid collision detected in batch generation")
)

// UUIDOptions defines the configuration for generating UUIDs.
type UUIDOptions struct {
	V7    bool
	Count int
}

// GenerateUUIDs produces one or more canonical UUIDs (v4 by default, or v7 with V7=true).
func GenerateUUIDs(opts UUIDOptions) ([]string, error) {
	if opts.Count <= 0 {
		return nil, ErrInvalidUUIDCount
	}

	results := make([]string, opts.Count)
	seen := make(map[string]struct{}, opts.Count)

	for i := 0; i < opts.Count; i++ {
		var u uuid.UUID
		var err error

		if opts.V7 {
			u, err = uuid.NewV7()
		} else {
			u, err = uuid.NewRandom()
		}

		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID: %w", err)
		}

		str := u.String()
		if _, exists := seen[str]; exists {
			return nil, ErrUUIDCollision
		}
		seen[str] = struct{}{}
		results[i] = str
	}

	return results, nil
}
