package util

import "strings"

// Contains checks if a slice contains a particular string (case-insensitive).
func Contains(slice []string, item string) bool {
    for _, s := range slice {
        if strings.EqualFold(s, item) {
            return true
        }
    }
    return false
}

// ParseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
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
