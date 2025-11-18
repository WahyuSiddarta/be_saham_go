package validator

import (
	"fmt"
	"time"
)

// parseDateString parses an optional date string in RFC3339 format into a time.Time pointer.
func parseDateString(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, *dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	return &t, nil
}
