// === pkg/util/util.go ===
// Package util provides utility functions used across the repo-to-txt tool.
// These functions include string manipulation and slice operations.
package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

// CopyFile copies a file from src to dst. If dst does not exist, it is created.
// If dst exists, it is overwritten.
//
// Parameters:
//   - src: Source file path.
//   - dst: Destination file path.
//
// Returns:
//   - error: An error if the copy fails.
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("unable to stat source file: %w", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("source file %s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open source file: %w", err)
	}
	defer source.Close()

	// Ensure the destination directory exists
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create destination directory: %w", err)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("error copying data: %w", err)
	}

	return nil
}
