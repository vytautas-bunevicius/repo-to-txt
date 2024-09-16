// Package clone_test contains unit tests for the clone package.
package clone

import (
	"testing"
)

// TestExtractRepoName verifies that the ExtractRepoName function correctly extracts
// repository names from various URL formats.
func TestExtractRepoName(t *testing.T) {
	testCases := []struct {
		repoURL  string
		expected string
		wantErr  bool
	}{
		{"https://github.com/user/repo.git", "repo", false},
		{"git@github.com:user/repo.git", "repo", false},
		{"invalid-url", "", true},
		{"ftp://github.com/user/repo.git", "", true},    // Invalid scheme
		{"https://github.com/user/repo", "repo", false}, // Without .git
		{"git@github.com:user/repo", "repo", false},     // SSH without .git
	}

	for _, tc := range testCases {
		t.Run(tc.repoURL, func(t *testing.T) {
			result, err := ExtractRepoName(tc.repoURL)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractRepoName(%q) error = %v, wantErr %v", tc.repoURL, err, tc.wantErr)
				return
			}
			if result != tc.expected {
				t.Errorf("ExtractRepoName(%q) = %q; want %q", tc.repoURL, result, tc.expected)
			}
		})
	}
}
