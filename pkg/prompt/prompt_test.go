package prompt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

// TestValidateRepoURL tests the validateRepoURL function
func TestValidateRepoURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid HTTPS URL", "https://github.com/user/repo.git", false},
		{"Valid SSH URL", "git@github.com:user/repo.git", false},
		{"Invalid URL", "http://github.com/user/repo.git", true},
		{"Empty URL", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRepoURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRepoURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIsSSHKeyPassphraseProtected tests the isSSHKeyPassphraseProtected function
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

// TestDefaultSSHKeyPath tests the defaultSSHKeyPath function
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

// TestDefaultDownloadsPath tests the defaultDownloadsPath function
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

// TestPromptForMissingInputs tests the PromptForMissingInputs function
func TestPromptForMissingInputs(t *testing.T) {
	// This is a basic test structure. In a real scenario, you'd need to mock the huh library,
	// which is beyond the scope of this example. Here's a simplified version:
	cfg := &config.Config{}
	err := PromptForMissingInputs(cfg)
	if err != nil {
		t.Fatalf("PromptForMissingInputs() error = %v", err)
	}

	// In a real test, you'd assert on the values in cfg here.
	// For now, we'll just check that the function ran without error.
	// You may want to expand this test when you implement proper mocking.
}
