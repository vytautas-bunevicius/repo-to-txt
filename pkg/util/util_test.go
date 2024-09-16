// Package util_test contains unit tests for the util package.
package util

import "testing"

// TestParseCommaSeparated verifies that the ParseCommaSeparated function correctly
// splits and trims comma-separated strings into a slice.
func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"a, b , c", []string{"a", "b", "c"}},
		{"  foo,bar , baz ", []string{"foo", "bar", "baz"}},
		{"", nil},
		{", ,", nil},
		{"single", []string{"single"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseCommaSeparated(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d items, got %d", len(tt.expected), len(result))
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("Expected %q at index %d, got %q", v, i, result[i])
				}
			}
		})
	}
}

// TestContains verifies that the Contains function correctly identifies
// the presence or absence of items in a slice.
func TestContains(t *testing.T) {
	tests := []struct {
		slice    []string
		item     string
		expected bool
	}{
		{[]string{"Go", "Python", "Java"}, "python", true},
		{[]string{"Go", "Python", "Java"}, "C++", false},
		{[]string{}, "any", false},
		{[]string{"Test"}, "test", true},
		{[]string{"CASE"}, "case", true},
		{[]string{"MixEd", "CaSe"}, "mIxEd", true},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Contains(%v, %q) = %v; want %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}
