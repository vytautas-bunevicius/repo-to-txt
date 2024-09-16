// Package config_test contains unit tests for the config package.
package config

import (
	"os"
	"testing"
)

// TestParseFlags verifies that the ParseFlags method correctly parses command-line flags
// and populates the Config struct accordingly.
func TestParseFlags(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		"cmd",
		"-repo=https://github.com/user/repo.git",
		"-auth=https",
		"-username=testuser",
		"-pat=testtoken",
		"-ssh-key=/path/to/ssh/key",
		"-output-dir=/path/to/output",
		"-exclude=docs,tests",
		"-include-ext=.go,.md",
		"-copy-clipboard=true",
	}

	cfg := NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		t.Fatalf("ParseFlags returned an error: %v", err)
	}

	if cfg.RepoURL != "https://github.com/user/repo.git" {
		t.Errorf("Expected RepoURL to be %q, got %q", "https://github.com/user/repo.git", cfg.RepoURL)
	}

	if cfg.AuthMethod != AuthMethodHTTPS {
		t.Errorf("Expected AuthMethod to be HTTPS, got %v", cfg.AuthMethod)
	}

	if cfg.Username != "testuser" {
		t.Errorf("Expected Username to be %q, got %q", "testuser", cfg.Username)
	}

	if cfg.PersonalAccessToken != "testtoken" {
		t.Errorf("Expected PersonalAccessToken to be %q, got %q", "testtoken", cfg.PersonalAccessToken)
	}

	if cfg.SSHKeyPath != "/path/to/ssh/key" {
		t.Errorf("Expected SSHKeyPath to be %q, got %q", "/path/to/ssh/key", cfg.SSHKeyPath)
	}

	if cfg.OutputDir != "/path/to/output" {
		t.Errorf("Expected OutputDir to be %q, got %q", "/path/to/output", cfg.OutputDir)
	}

	expectedExcludes := []string{"docs", "tests"}
	if len(cfg.ExcludeFolders) != len(expectedExcludes) {
		t.Errorf("Expected ExcludeFolders length to be %d, got %d", len(expectedExcludes), len(cfg.ExcludeFolders))
	} else {
		for i, v := range expectedExcludes {
			if cfg.ExcludeFolders[i] != v {
				t.Errorf("Expected ExcludeFolders[%d] to be %q, got %q", i, v, cfg.ExcludeFolders[i])
			}
		}
	}

	expectedIncludes := []string{".go", ".md"}
	if len(cfg.IncludeExt) != len(expectedIncludes) {
		t.Errorf("Expected IncludeExt length to be %d, got %d", len(expectedIncludes), len(cfg.IncludeExt))
	} else {
		for i, v := range expectedIncludes {
			if cfg.IncludeExt[i] != v {
				t.Errorf("Expected IncludeExt[%d] to be %q, got %q", i, v, cfg.IncludeExt[i])
			}
		}
	}

	if !cfg.CopyToClipboard {
		t.Errorf("Expected CopyToClipboard to be true, got false")
	}
}

// TestParseFlagsDefaults verifies that default values are correctly set when certain flags are omitted.
func TestParseFlagsDefaults(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		"cmd",
		"-repo=https://github.com/user/repo.git",
	}

	cfg := NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		t.Fatalf("ParseFlags returned an error: %v", err)
	}

	if cfg.RepoURL != "https://github.com/user/repo.git" {
		t.Errorf("Expected RepoURL to be %q, got %q", "https://github.com/user/repo.git", cfg.RepoURL)
	}

	if cfg.AuthMethod != AuthMethodNone {
		t.Errorf("Expected AuthMethod to be None, got %v", cfg.AuthMethod)
	}

	if cfg.OutputDir != "" {
		t.Errorf("Expected OutputDir to be empty, got %q", cfg.OutputDir)
	}

	if len(cfg.ExcludeFolders) != 0 {
		t.Errorf("Expected ExcludeFolders to be empty, got %v", cfg.ExcludeFolders)
	}

	if len(cfg.IncludeExt) != 0 {
		t.Errorf("Expected IncludeExt to be empty, got %v", cfg.IncludeExt)
	}

	if cfg.CopyToClipboard {
		t.Errorf("Expected CopyToClipboard to be false, got true")
	}
}

// TestParseFlagsInvalidAuth verifies that an invalid authentication method results in an error.
func TestParseFlagsInvalidAuth(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		"cmd",
		"-repo=https://github.com/user/repo.git",
		"-auth=invalid",
	}

	cfg := NewConfig()
	err := cfg.ParseFlags()
	if err == nil {
		t.Errorf("Expected ParseFlags to return an error for invalid auth method, got nil")
	}
}

