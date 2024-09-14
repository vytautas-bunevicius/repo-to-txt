package config

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd", "-repo=https://github.com/user/repo.git", "-auth=https"}

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
}
