// Package config handles the configuration for the repo-to-txt CLI tool.
// It parses command-line flags and manages authentication and output settings.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	// Version is the current version of the tool.
	Version = "1.1.0"

	// DefaultCloneDir is the default directory name for cloning.
	DefaultCloneDir = "repo-to-txt-clone"

	// DefaultOutputExt is the default extension for the output file.
	DefaultOutputExt = ".txt"

	// DefaultSSHKeyName is the default SSH key name.
	DefaultSSHKeyName = "git"

	// DefaultExcludedExt is the default file extension to exclude.
	DefaultExcludedExt = ".ipynb"
)

// AuthMethod represents the type of authentication to use.
type AuthMethod int

const (
	// AuthMethodNone indicates no authentication.
	AuthMethodNone AuthMethod = iota

	// AuthMethodHTTPS indicates HTTPS authentication.
	AuthMethodHTTPS

	// AuthMethodSSH indicates SSH authentication.
	AuthMethodSSH
)

// Config holds the configuration for the CLI tool.
type Config struct {
	RepoURL             string
	AuthMethod          AuthMethod
	Username            string
	PersonalAccessToken string
	SSHKeyPath          string
	SSHPassphrase       string
	ExcludeFolders      []string
	IncludeExt          []string
	OutputDir           string
	CopyToClipboard     bool // New field for clipboard copying
	AuthFlagSet         bool
	VersionFlag         bool
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{}
}

// ParseFlags parses command-line flags and populates the Config struct.
// It handles required flags, default values, and validates authentication methods.
//
// Returns:
//   - error: An error if flag parsing or validation fails.
func (cfg *Config) ParseFlags() error {
	var authMethod string
	var excludeFolders, includeExt string
	var copyToClipboard bool // New variable for clipboard flag

	flag.StringVar(&cfg.RepoURL, "repo", "", "GitHub repository URL (HTTPS or SSH) (Required)")
	flag.StringVar(&authMethod, "auth", "", "Authentication method: none, https, or ssh (Required)")
	flag.StringVar(&cfg.Username, "username", "", "GitHub username (for HTTPS)")
	flag.StringVar(&cfg.PersonalAccessToken, "pat", "", "GitHub Personal Access Token (for HTTPS)")
	flag.StringVar(&cfg.SSHKeyPath, "ssh-key", "", "Path to SSH private key (for SSH)")
	flag.StringVar(&cfg.OutputDir, "output-dir", "", "Output directory for the generated text file")
	flag.StringVar(&excludeFolders, "exclude", "", "Comma-separated list of folders to exclude from the output")
	flag.StringVar(&includeExt, "include-ext", "", "Comma-separated list of file extensions to include (e.g., .go,.md). If not set, defaults to excluding certain non-code files like .ipynb")
	flag.BoolVar(&copyToClipboard, "copy-clipboard", false, "Copy the output to the clipboard after creation") // New flag
	flag.BoolVar(&cfg.VersionFlag, "version", false, "Print the version number and exit")

	flag.Parse()

	if cfg.VersionFlag {
		fmt.Printf("repo-to-txt version %s\n", Version)
		os.Exit(0) // Exit after printing version
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "auth" {
			cfg.AuthFlagSet = true
		}
	})

	cfg.ExcludeFolders = parseCommaSeparated(excludeFolders)
	cfg.IncludeExt = parseCommaSeparated(includeExt)
	cfg.CopyToClipboard = copyToClipboard // Assign the flag value

	switch strings.ToLower(authMethod) {
	case "https":
		cfg.AuthMethod = AuthMethodHTTPS
	case "ssh":
		cfg.AuthMethod = AuthMethodSSH
	case "none", "":
		cfg.AuthMethod = AuthMethodNone
	default:
		return errors.New("invalid authentication method: choose from none, https, ssh")
	}

	return nil
}

// parseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
//
// Parameters:
//   - input: The comma-separated string.
//
// Returns:
//   - []string: A slice of trimmed strings.
func parseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
