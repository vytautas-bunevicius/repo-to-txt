package prompt

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/util"
	"golang.org/x/term"
)

// ErrEmptyInput is returned when the user provides an empty input.
var ErrEmptyInput = errors.New("input cannot be empty")

// PromptForMissingInputs prompts the user for any missing configuration inputs.
func PromptForMissingInputs(cfg *config.Config) error {
	if cfg.RepoURL == "" {
		cfg.RepoURL = promptForRepoURL()
	}

	if !cfg.AuthFlagSet {
		cfg.AuthMethod = promptForAuthMethod()
	}

	switch cfg.AuthMethod {
	case config.AuthMethodHTTPS:
		if cfg.Username == "" {
			cfg.Username = promptForInput("Enter your GitHub username:")
		}
		if cfg.PersonalAccessToken == "" {
			cfg.PersonalAccessToken = promptForPassword("Enter your GitHub Personal Access Token:")
		}
	case config.AuthMethodSSH:
		if cfg.SSHKeyPath == "" {
			cfg.SSHKeyPath = promptForInputWithDefault("Enter the path to your SSH private key", defaultSSHKeyPath())
		}
		if isSSHKeyPassphraseProtected(cfg.SSHKeyPath) {
			cfg.SSHPassphrase = promptForPassword("Enter your SSH key passphrase (leave empty if none):")
		}
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = promptForInputWithDefault("Enter the output directory", ".")
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(cfg.OutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if len(cfg.ExcludeFolders) == 0 {
		excludeInput := promptForInputWithDefault("Enter folders to exclude (comma-separated, leave empty to include all)", "")
		cfg.ExcludeFolders = util.ParseCommaSeparated(excludeInput)
	}

	if len(cfg.IncludeExt) == 0 {
		includeInput := promptForInputWithDefault("Enter file extensions to include (comma-separated, leave empty to include all)", "")
		cfg.IncludeExt = util.ParseCommaSeparated(includeInput)
	}

	return validateRepoURL(cfg.RepoURL)
}

// promptForRepoURL prompts the user for a valid repository URL.
func promptForRepoURL() string {
	for {
		input := promptForInput("Enter the GitHub repository URL (HTTPS or SSH):")
		if err := validateRepoURL(input); err == nil {
			return input
		}
		fmt.Println("Invalid repository URL. Please enter a valid HTTPS or SSH URL for a GitHub repository.")
	}
}

// promptForInput prompts the user for input with the given label.
func promptForInput(label string) string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validateNonEmpty,
		Stdin:    os.Stdin,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v", err)
	}
	return strings.TrimSpace(result)
}

// promptForInputWithDefault prompts the user for input with a default value.
func promptForInputWithDefault(label, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
		Stdin:   os.Stdin,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v", err)
	}
	input := strings.TrimSpace(result)
	if input == "" {
		return defaultValue
	}
	return input
}

// promptForAuthMethod prompts the user to select an authentication method.
func promptForAuthMethod() config.AuthMethod {
	prompt := promptui.Select{
		Label: "Select authentication method",
		Items: []string{"No Authentication", "HTTPS with PAT", "SSH"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v", err)
	}

	switch result {
	case "HTTPS with PAT":
		return config.AuthMethodHTTPS
	case "SSH":
		return config.AuthMethodSSH
	default:
		return config.AuthMethodNone
	}
}

// promptForPassword securely prompts the user for a password without echoing.
func promptForPassword(label string) string {
	fmt.Printf("%s ", label)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Password prompt failed: %v", err)
	}
	fmt.Println() // Move to the next line after input
	return strings.TrimSpace(string(bytePassword))
}

// validateNonEmpty ensures the input is not empty.
func validateNonEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return ErrEmptyInput
	}
	return nil
}

// defaultSSHKeyPath returns the default SSH key path.
func defaultSSHKeyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to determine user home directory: %v", err)
	}
	return filepath.Join(home, ".ssh", "id_rsa")
}

// isSSHKeyPassphraseProtected checks if the SSH key is passphrase protected.
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

// validateRepoURL validates the repository URL.
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
