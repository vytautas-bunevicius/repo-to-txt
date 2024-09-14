// Package util provides utility functions used across the repo-to-txt tool.
// These functions include string manipulation and slice operations.
package util

import "strings"

// Contains checks if a slice contains a particular string (case-insensitive).
//
// Parameters:
//   - slice: The slice of strings to search within.
//   - item: The string to search for.
//
// Returns:
//   - bool: True if the slice contains the item, false otherwise.
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// ParseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
//
// Parameters:
//   - input: The comma-separated string.
//
// Returns:
//   - []string: A slice of trimmed strings. Returns nil if the input is empty.
func ParseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
