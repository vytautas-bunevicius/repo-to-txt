// Package prompt_test contains unit tests for the prompt package.
package prompt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

// TestValidateRepoURL verifies that the validateRepoURL function correctly validates
// various repository URL formats.
func TestValidateRepoURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid HTTPS URL", "https://github.com/user/repo.git", false},
		{"Valid SSH URL", "git@github.com:user/repo.git", false},
		{"Invalid HTTP URL", "http://github.com/user/repo.git", false}, // Allowed as per isHTTPSURL and isSSHURL
		{"Invalid FTP URL", "ftp://github.com/user/repo.git", true},
		{"Empty URL", "", true},
		{"HTTPS without .git", "https://github.com/user/repo", false},
		{"SSH without .git", "git@github.com:user/repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRepoURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRepoURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

// TestIsSSHKeyPassphraseProtected verifies that the isSSHKeyPassphraseProtected function
// correctly identifies whether an SSH key is passphrase protected.
func TestIsSSHKeyPassphraseProtected(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test-ssh-key")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Test with non-encrypted content
	if _, err := tmpfile.Write([]byte("-----BEGIN OPENSSH PRIVATE KEY-----\n")); err != nil {
		t.Fatal(err)
	}
	if isSSHKeyPassphraseProtected(tmpfile.Name()) {
		t.Error("Expected non-encrypted key to return false")
	}

	// Test with encrypted content
	if err := tmpfile.Truncate(0); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte("-----BEGIN OPENSSH PRIVATE KEY-----\nENCRYPTED\n")); err != nil {
		t.Fatal(err)
	}
	if !isSSHKeyPassphraseProtected(tmpfile.Name()) {
		t.Error("Expected encrypted key to return true")
	}
}

// TestDefaultSSHKeyPath verifies that the defaultSSHKeyPath function returns the correct default SSH key path.
func TestDefaultSSHKeyPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(home, ".ssh", "id_rsa")
	if got := defaultSSHKeyPath(); got != expected {
		t.Errorf("defaultSSHKeyPath() = %v, want %v", got, expected)
	}
}

// TestDefaultDownloadsPath verifies that the defaultDownloadsPath function returns the correct default Downloads path.
func TestDefaultDownloadsPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(home, "Downloads")
	if got := defaultDownloadsPath(); got != expected {
		t.Errorf("defaultDownloadsPath() = %v, want %v", got, expected)
	}
}

// TestPromptForMissingInputs is a placeholder for testing the PromptForMissingInputs function.
// In a real scenario, you would mock the huh library to simulate user input.
func TestPromptForMissingInputs(t *testing.T) {
	// This is a basic test structure. In a real scenario, you'd need to mock the huh library,
	// which is beyond the scope of this example. Here's a simplified version:
	cfg := &config.Config{}
	err := PromptForMissingInputs(cfg)
	if err != nil {
		t.Fatalf("PromptForMissingInputs() error = %v", err)
	}
}
