package util

import "testing"

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
