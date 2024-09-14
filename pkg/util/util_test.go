// Package util_test contains unit tests for the util package.
package util

import "testing"

// TestParseCommaSeparated verifies that the ParseCommaSeparated function correctly
// splits and trims comma-separated strings into a slice.
func TestParseCommaSeparated(t *testing.T) {
	input := "a, b , c"
	expected := []string{"a", "b", "c"}

	result := ParseCommaSeparated(input)
	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %q at index %d, got %q", v, i, result[i])
		}
	}
}
