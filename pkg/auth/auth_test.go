// Package auth_test contains unit tests for the auth package.
package auth

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh" // Added import for ssh package
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

// TestSetupAuth verifies that the SetupAuth function correctly sets up the authentication method
// based on different configuration scenarios.
func TestSetupAuth(t *testing.T) {
	// Test HTTPS Authentication
	cfgHTTPS := &config.Config{
		AuthMethod:          config.AuthMethodHTTPS,
		Username:            "testuser",
		PersonalAccessToken: "testtoken",
	}
	authMethodHTTPS, err := SetupAuth(cfgHTTPS)
	if err != nil {
		t.Fatalf("SetupAuth returned an error: %v", err)
	}

	if authMethodHTTPS == nil {
		t.Errorf("Expected HTTPS AuthMethod, got nil")
	}

	// Assert that authMethodHTTPS is of type *ssh.PublicKeys
	// Note: In go-git, HTTPS uses http.BasicAuth, not ssh.PublicKeys.
	// Correcting the expectation.
	if _, ok := authMethodHTTPS.(*ssh.PublicKeys); ok {
		t.Errorf("Expected HTTPS AuthMethod to be of type *http.BasicAuth, got *ssh.PublicKeys")
	}

	// Test SSH Authentication without passphrase
	cfgSSH := &config.Config{
		AuthMethod: config.AuthMethodSSH,
		SSHKeyPath: "/path/to/ssh/key",
	}
	authMethodSSH, err := SetupAuth(cfgSSH)
	if err != nil {
		t.Fatalf("SetupAuth returned an error: %v", err)
	}

	if authMethodSSH == nil {
		t.Errorf("Expected SSH AuthMethod, got nil")
	}

	// Assert that authMethodSSH is of type *ssh.PublicKeys
	if _, ok := authMethodSSH.(*ssh.PublicKeys); !ok {
		t.Errorf("Expected SSH AuthMethod to be of type *ssh.PublicKeys, got %T", authMethodSSH)
	}

	// Test SSH Authentication with passphrase
	cfgSSHPass := &config.Config{
		AuthMethod:    config.AuthMethodSSH,
		SSHKeyPath:    "/path/to/ssh/key",
		SSHPassphrase: "passphrase",
	}
	authMethodSSHPass, err := SetupAuth(cfgSSHPass)
	if err != nil {
		t.Fatalf("SetupAuth returned an error: %v", err)
	}

	if authMethodSSHPass == nil {
		t.Errorf("Expected SSH AuthMethod with passphrase, got nil")
	}

	// Assert that authMethodSSHPass is of type *ssh.PublicKeys
	if _, ok := authMethodSSHPass.(*ssh.PublicKeys); !ok {
		t.Errorf("Expected SSH AuthMethod to be of type *ssh.PublicKeys, got %T", authMethodSSHPass)
	}

	// Test No Authentication
	cfgNone := &config.Config{
		AuthMethod: config.AuthMethodNone,
	}
	authMethodNone, err := SetupAuth(cfgNone)
	if err != nil {
		t.Fatalf("SetupAuth returned an error: %v", err)
	}

	if authMethodNone != nil {
		t.Errorf("Expected No AuthMethod to be nil, got %T", authMethodNone)
	}

	// Test Missing HTTPS Credentials
	cfgMissingHTTPS := &config.Config{
		AuthMethod: config.AuthMethodHTTPS,
	}
	_, err = SetupAuth(cfgMissingHTTPS)
	if err == nil {
		t.Errorf("Expected error for missing HTTPS credentials, got nil")
	}

	// Test Unsupported Authentication Method
	cfgInvalid := &config.Config{
		AuthMethod: 999, // Invalid AuthMethod
	}
	_, err = SetupAuth(cfgInvalid)
	if err == nil {
		t.Errorf("Expected error for unsupported authentication method, got nil")
	}
}
