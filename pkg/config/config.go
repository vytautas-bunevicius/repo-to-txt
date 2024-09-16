// Package config handles the configuration for the repo-to-txt CLI tool.
// It defines the structure for storing configuration options and provides
// methods for parsing command-line flags and environment variables.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Constants define default values and version information for the tool.
const (
	// Version represents the current version of the tool.
	Version = "1.1.0"

	// DefaultCloneDir is the default directory name for cloning repositories.
	DefaultCloneDir = "repo-to-txt-clone"

	// DefaultOutputExt is the default file extension for the output file.
	DefaultOutputExt = ".txt"

	// DefaultSSHKeyName is the default name for the SSH key file.
	DefaultSSHKeyName = "git"

	// DefaultExcludedExt is the default file extension to exclude from processing.
	DefaultExcludedExt = ".ipynb"
)

// AuthMethod represents the type of authentication to use when accessing repositories.
type AuthMethod int

// Constants for different authentication methods.
const (
	AuthMethodNone AuthMethod = iota
	AuthMethodHTTPS
	AuthMethodSSH
)

// Config holds all configuration options for the repo-to-txt tool.
type Config struct {
	RepoURL             string     // URL of the Git repository to clone
	AuthMethod          AuthMethod // Authentication method to use
	Username            string     // GitHub username for HTTPS authentication
	PersonalAccessToken string     // GitHub personal access token for HTTPS authentication
	SSHKeyPath          string     // Path to SSH key for SSH authentication
	SSHPassphrase       string     // Passphrase for SSH key, if any
	ExcludeFolders      []string   // List of folders to exclude from processing
	IncludeExt          []string   // List of file extensions to include in processing
	FileNames           []string   // List of exact file names to copy from the repository
	OutputDir           string     // Directory to output the generated text file
	AuthFlagSet         bool       // Indicates if authentication method was set via flag
	VersionFlag         bool       // Flag to print version information
	CopyToClipboard     bool       // Flag to copy output to clipboard
	CopyToClipboardSet  bool       // Indicates if copy-to-clipboard was set via flag
}

// NewConfig creates and returns a new Config instance with default values.
func NewConfig() *Config {
	return &Config{}
}

// ParseFlags parses command-line flags and populates the Config struct.
// It handles required flags, default values, and validates authentication methods.
func (cfg *Config) ParseFlags() error {
	var authMethod string
	var excludeFolders, includeExt, files string

	// Define command-line flags
	flag.StringVar(&cfg.RepoURL, "repo", "", "GitHub repository URL (HTTPS or SSH) (Required)")
	flag.StringVar(&authMethod, "auth", "", "Authentication method: none, https, or ssh (Required)")
	flag.StringVar(&cfg.Username, "username", "", "GitHub username (for HTTPS)")
	flag.StringVar(&cfg.PersonalAccessToken, "pat", "", "GitHub Personal Access Token (for HTTPS)")
	flag.StringVar(&cfg.SSHKeyPath, "ssh-key", "", "Path to SSH private key (for SSH)")
	flag.StringVar(&cfg.OutputDir, "output-dir", "", "Output directory for the generated text file")
	flag.StringVar(&excludeFolders, "exclude", "", "Comma-separated list of folders to exclude from the output")
	flag.StringVar(&includeExt, "include-ext", "", "Comma-separated list of file extensions to include (e.g., .go,.md). If not set, defaults to excluding certain non-code files like .ipynb")
	flag.StringVar(&files, "files", "", "Comma-separated list of exact file names to copy from the repository")
	flag.BoolVar(&cfg.VersionFlag, "version", false, "Print the version number and exit")
	flag.BoolVar(&cfg.CopyToClipboard, "copy-clipboard", false, "Copy the output to clipboard")

	// Parse the flags
	flag.Parse()

	// Check if copy-to-clipboard was set via flag
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "copy-clipboard" {
			cfg.CopyToClipboardSet = true
		}
		if f.Name == "auth" {
			cfg.AuthFlagSet = true
		}
	})

	// Handle version flag
	if cfg.VersionFlag {
		fmt.Printf("repo-to-txt version %s\n", Version)
		os.Exit(0)
	}

	// Process comma-separated inputs
	cfg.ExcludeFolders = parseCommaSeparated(excludeFolders)
	cfg.IncludeExt = parseCommaSeparated(includeExt)
	cfg.FileNames = parseCommaSeparated(files)

	// Set authentication method
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
// It returns nil if the input string is empty.
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
