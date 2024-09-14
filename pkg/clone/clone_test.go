package clone

import (
	"testing"
)

func TestExtractRepoName(t *testing.T) {
	testCases := []struct {
		repoURL  string
		expected string
		wantErr  bool
	}{
		{"https://github.com/user/repo.git", "repo", false},
		{"git@github.com:user/repo.git", "repo", false},
		{"invalid-url", "", true},
		{"ftp://github.com/user/repo.git", "", true}, // Invalid scheme
	}

	for _, tc := range testCases {
		result, err := ExtractRepoName(tc.repoURL)
		if (err != nil) != tc.wantErr {
			t.Errorf("ExtractRepoName(%q) error = %v, wantErr %v", tc.repoURL, err, tc.wantErr)
			continue
		}
		if result != tc.expected {
			t.Errorf("ExtractRepoName(%q) = %q; want %q", tc.repoURL, result, tc.expected)
		}
	}
}
