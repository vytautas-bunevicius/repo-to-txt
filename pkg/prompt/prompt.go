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
	// Step 1: Prompt for the GitHub repository URL.
	repoForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("GitHub repository URL (HTTPS or SSH)").
				Value(&cfg.RepoURL).
				Validate(validateRepoURL),
		),
	)
	err := repoForm.Run()
	if err != nil {
		return fmt.Errorf("repository URL input error: %w", err)
	}

	// Determine the authentication options based on the repository URL type.
	var authOptions []huh.Option[config.AuthMethod]
	if isHTTPSURL(cfg.RepoURL) {
		authOptions = []huh.Option[config.AuthMethod]{
			huh.NewOption("No Authentication", config.AuthMethodNone),
			huh.NewOption("HTTPS with PAT", config.AuthMethodHTTPS),
		}
	} else if isSSHURL(cfg.RepoURL) {
		authOptions = []huh.Option[config.AuthMethod]{
			huh.NewOption("SSH Authentication", config.AuthMethodSSH),
		}
	} else {
		return errors.New("unsupported repository URL format for authentication options")
	}

	// Step 2: Prompt for the authentication method based on the repository URL type.
	authForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[config.AuthMethod]().
				Title("Select authentication method").
				Options(authOptions...).
				Value(&cfg.AuthMethod),
		),
	)
	err = authForm.Run()
	if err != nil {
		return fmt.Errorf("authentication method input error: %w", err)
	}

	// Based on the selected authentication method, prompt for additional inputs.
	switch cfg.AuthMethod {
	case config.AuthMethodHTTPS:
		httpsForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("GitHub username").
					Value(&cfg.Username).
					Validate(nonEmptyValidator("GitHub username")),
				huh.NewInput().
					Title("GitHub Personal Access Token").
					Value(&cfg.PersonalAccessToken).
					Password(true).
					Validate(nonEmptyValidator("Personal Access Token")),
			),
		)
		err = httpsForm.Run()
		if err != nil {
			return fmt.Errorf("HTTPS authentication input error: %w", err)
		}
	case config.AuthMethodSSH:
		defaultSSHKey := defaultSSHKeyPath()
		sshForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Path to SSH private key").
					Value(&cfg.SSHKeyPath).
					Placeholder(defaultSSHKey).
					Validate(func(s string) error {
						if s == "" {
							cfg.SSHKeyPath = defaultSSHKey
							return nil
						}
						if _, err := os.Stat(s); os.IsNotExist(err) {
							return fmt.Errorf("SSH key file does not exist at path: %s", s)
						}
						return nil
					}),
			),
		)
		err = sshForm.Run()
		if err != nil {
			return fmt.Errorf("SSH key path input error: %w", err)
		}

		// Check if the SSH key is passphrase protected.
		if isSSHKeyPassphraseProtected(cfg.SSHKeyPath) {
			passphraseForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("SSH key passphrase (leave empty if none)").
						Value(&cfg.SSHPassphrase).
						Password(true),
				),
			)
			err = passphraseForm.Run()
			if err != nil {
				return fmt.Errorf("SSH passphrase input error: %w", err)
			}
		}
	}

	// Step 3: Prompt for output configurations.
	var excludeFolders, includeExt string
	defaultOutputDir := defaultDownloadsPath()
	outputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Output directory").
				Value(&cfg.OutputDir).
				Placeholder(defaultOutputDir).
				Validate(func(s string) error {
					if s == "" {
						cfg.OutputDir = defaultOutputDir
						return nil
					}
					return nil
				}),
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
		return fmt.Errorf("output configuration input error: %w", err)
	}

	// Process the comma-separated inputs.
	cfg.ExcludeFolders = util.ParseCommaSeparated(excludeFolders)
	cfg.IncludeExt = util.ParseCommaSeparated(includeExt)

	// Ensure the output directory exists.
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

	if isHTTPSURL(repoURL) || isSSHURL(repoURL) {
		return nil
	}

	return errors.New("URL must be either HTTPS (https://github.com/user/repo) or SSH (git@github.com:user/repo) format")
}

// isHTTPSURL checks if the provided URL is an HTTPS GitHub repository URL.
//
// Parameters:
//   - url: The repository URL to check.
//
// Returns:
//   - bool: True if the URL starts with "https://github.com/", false otherwise.
func isHTTPSURL(url string) bool {
	return strings.HasPrefix(url, "https://github.com/")
}

// isSSHURL checks if the provided URL is an SSH GitHub repository URL.
//
// Parameters:
//   - url: The repository URL to check.
//
// Returns:
//   - bool: True if the URL starts with "git@github.com:", false otherwise.
func isSSHURL(url string) bool {
	return strings.HasPrefix(url, "git@github.com:")
}

// nonEmptyValidator returns a validator function that ensures the input string is not empty.
//
// Parameters:
//   - fieldName: The name of the field being validated.
//
// Returns:
//   - func(string) error: A validator function.
func nonEmptyValidator(fieldName string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s cannot be empty", fieldName)
		}
		return nil
	}
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

	// Attempt to read the first few bytes to check for encryption.
	buf := make([]byte, 100)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	content := string(buf[:n])
	return strings.Contains(content, "ENCRYPTED")
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
