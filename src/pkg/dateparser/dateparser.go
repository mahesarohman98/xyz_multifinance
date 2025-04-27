package dateparser

import (
	"fmt"
	"time"
)

// ParseDate parses a date string with a 2006-01-02 layout
func ParseDate(dateStr string) (time.Time, error) {
	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %w", err)
	}
	return date, nil
}
