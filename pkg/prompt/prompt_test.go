package prompt

import (
	"testing"
)

func TestValidateNonEmpty(t *testing.T) {
	if err := validateNonEmpty("test"); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := validateNonEmpty("   "); err == nil {
		t.Errorf("Expected an error for empty input, got nil")
	}
}
