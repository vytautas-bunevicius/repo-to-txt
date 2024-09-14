// Package prompt manages interactive user prompts for missing configuration inputs.
// It utilizes the huh library to create forms for collecting user input.
package prompt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/util"
)

// ErrEmptyInput is returned when the user provides an empty input.
var ErrEmptyInput = errors.New("input cannot be empty")

// PromptForMissingInputs prompts the user interactively for any missing configuration inputs.
// It updates the provided Config struct with the collected inputs.
//
// Parameters:
//   - cfg: A pointer to the Config struct to be populated.
//
// Returns:
//   - error: An error if prompting fails or input validation fails.
func PromptForMissingInputs(cfg *config.Config) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("GitHub repository URL (HTTPS or SSH)").
				Value(&cfg.RepoURL).
				Validate(validateRepoURL),

			huh.NewSelect[config.AuthMethod]().
				Title("Select authentication method").
				Options(
					huh.NewOption("No Authentication", config.AuthMethodNone),
					huh.NewOption("HTTPS with PAT", config.AuthMethodHTTPS),
					huh.NewOption("SSH", config.AuthMethodSSH),
				).
				Value(&cfg.AuthMethod),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("form input error: %w", err)
	}

	switch cfg.AuthMethod {
	case config.AuthMethodHTTPS:
		httpsForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("GitHub username").
					Value(&cfg.Username),
				huh.NewInput().
					Title("GitHub Personal Access Token").
					Value(&cfg.PersonalAccessToken).
					Password(true),
			),
		)
		err = httpsForm.Run()
	case config.AuthMethodSSH:
		defaultSSHKeyPath := defaultSSHKeyPath()
		sshForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Path to SSH private key").
					Value(&cfg.SSHKeyPath).
					Placeholder(defaultSSHKeyPath).
					Validate(func(s string) error {
						if s == "" {
							cfg.SSHKeyPath = defaultSSHKeyPath
							return nil
						}
						return nil
					}),
			),
		)
		err = sshForm.Run()
		if err == nil && isSSHKeyPassphraseProtected(cfg.SSHKeyPath) {
			passphraseForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("SSH key passphrase (leave empty if none)").
						Value(&cfg.SSHPassphrase).
						Password(true),
				),
			)
			err = passphraseForm.Run()
		}
	}

	if err != nil {
		return fmt.Errorf("authentication input error: %w", err)
	}

	var excludeFolders, includeExt string
	defaultOutputDir := defaultDownloadsPath()
	outputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Output directory").
				Value(&cfg.OutputDir).
				Placeholder(defaultOutputDir),
			huh.NewInput().
				Title("Folders to exclude (comma-separated, leave empty to include all)").
				Value(&excludeFolders),
			huh.NewInput().
				Title("File extensions to include (comma-separated, leave empty to include all)").
				Value(&includeExt),
		),
	)

	err = outputForm.Run()
	if err != nil {
		return fmt.Errorf("output configuration error: %w", err)
	}

	// Process the comma-separated inputs
	cfg.ExcludeFolders = util.ParseCommaSeparated(excludeFolders)
	cfg.IncludeExt = util.ParseCommaSeparated(includeExt)

	// Set default output directory if not provided
	if cfg.OutputDir == "" {
		cfg.OutputDir = defaultOutputDir
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(cfg.OutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}

// validateRepoURL validates the format of the provided GitHub repository URL.
// It ensures that the URL is either a valid HTTPS or SSH GitHub repository URL.
//
// Parameters:
//   - repoURL: The repository URL to validate.
//
// Returns:
//   - error: An error if the URL is invalid.
func validateRepoURL(repoURL string) error {
	if strings.TrimSpace(repoURL) == "" {
		return errors.New("repository URL cannot be empty")
	}

	if strings.HasPrefix(repoURL, "https://") {
		if !strings.HasPrefix(repoURL, "https://github.com/") {
			return errors.New("HTTPS URL must be a GitHub repository URL")
		}
		return nil
	}

	if strings.HasPrefix(repoURL, "git@github.com:") {
		return nil
	}

	return errors.New("URL must be either HTTPS (https://github.com/user/repo) or SSH (git@github.com:user/repo) format")
}

// isSSHKeyPassphraseProtected checks if the SSH key at the given path is protected by a passphrase.
// It does this by looking for the "ENCRYPTED" keyword in the key file.
//
// Parameters:
//   - keyPath: The file system path to the SSH private key.
//
// Returns:
//   - bool: True if the key is passphrase protected, false otherwise.
func isSSHKeyPassphraseProtected(keyPath string) bool {
	file, err := os.Open(keyPath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Attempt to read the first few bytes to check for encryption
	buf := make([]byte, 100)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	content := string(buf[:n])
	return strings.Contains(content, "ENCRYPTED")
}

// defaultSSHKeyPath returns the default path to the SSH private key.
// It typically points to the user's home directory under .ssh/id_rsa.
//
// Returns:
//   - string: The default SSH key path.
func defaultSSHKeyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".ssh", "id_rsa")
}

// defaultDownloadsPath returns the default path to the Downloads directory in the user's home.
// If the home directory cannot be determined, it defaults to the current directory.
//
// Returns:
//   - string: The default Downloads path.
func defaultDownloadsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, "Downloads")
}
