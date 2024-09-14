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

// Authentication methods.
const (
	AuthMethodNone AuthMethod = iota
	AuthMethodHTTPS
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
	AuthFlagSet         bool
	VersionFlag         bool
}

// NewConfig creates a new Config instance.
func NewConfig() *Config {
	return &Config{}
}

// ParseFlags parses command-line flags into the Config struct.
func (cfg *Config) ParseFlags() error {
	var authMethod string
	var excludeFolders, includeExt string

	flag.StringVar(&cfg.RepoURL, "repo", "", "GitHub repository URL (HTTPS or SSH) (Required)")
	flag.StringVar(&authMethod, "auth", "", "Authentication method: none, https, or ssh (Required)")
	flag.StringVar(&cfg.Username, "username", "", "GitHub username (for HTTPS)")
	flag.StringVar(&cfg.PersonalAccessToken, "pat", "", "GitHub Personal Access Token (for HTTPS)")
	flag.StringVar(&cfg.SSHKeyPath, "ssh-key", "", "Path to SSH private key (for SSH)")
	flag.StringVar(&cfg.OutputDir, "output-dir", "", "Output directory for the generated text file")
	flag.StringVar(&excludeFolders, "exclude", "", "Comma-separated list of folders to exclude from the output")
	flag.StringVar(&includeExt, "include-ext", "", "Comma-separated list of file extensions to include (e.g., .go,.md). If not set, defaults to excluding certain non-code files like .ipynb")
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

// Helper function to parse comma-separated values.
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
