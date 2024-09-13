// Package main provides a CLI tool to clone a GitHub repository
// and write its contents to a text file.
package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

// Version is the current version of the tool.
const Version = "1.0.1"

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
}

// AuthMethod represents the type of authentication to use.
type AuthMethod int

// Constants for different authentication methods and default values.
const (
	AuthMethodNone AuthMethod = iota
	AuthMethodHTTPS
	AuthMethodSSH

	defaultCloneDir    = "repo-to-txt-clone"
	defaultOutputExt   = ".txt"
	defaultSSHKeyName  = "git"
	defaultExcludedExt = ".ipynb" // Default exclusion for Jupyter notebooks
)

var (
	config      Config
	authFlagSet bool

	// ErrEmptyInput is returned when the user provides an empty input.
	ErrEmptyInput = errors.New("input cannot be empty")
)

func init() {
	flag.StringVar(&config.RepoURL, "repo", "", "GitHub repository URL (HTTPS or SSH) (Required)")
	authMethod := flag.String("auth", "", "Authentication method: none, https, or ssh (Required)")
	flag.StringVar(&config.Username, "username", "", "GitHub username (for HTTPS)")
	flag.StringVar(&config.PersonalAccessToken, "pat", "", "GitHub Personal Access Token (for HTTPS)")
	flag.StringVar(&config.SSHKeyPath, "ssh-key", "", "Path to SSH private key (for SSH)")
	var excludeFolders, includeExt string
	flag.StringVar(&excludeFolders, "exclude", "", "Comma-separated list of folders to exclude from the output")
	flag.StringVar(&includeExt, "include-ext", "", "Comma-separated list of file extensions to include (e.g., .go,.md). If not set, defaults to excluding certain non-code files like .ipynb")
	versionFlag := flag.Bool("version", false, "Print the version number and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("repo-to-txt version %s\n", Version)
		os.Exit(0)
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "auth" {
			authFlagSet = true
		}
	})

	config.ExcludeFolders = parseCommaSeparated(excludeFolders)
	config.IncludeExt = parseCommaSeparated(includeExt)

	switch strings.ToLower(*authMethod) {
	case "https":
		config.AuthMethod = AuthMethodHTTPS
	case "ssh":
		config.AuthMethod = AuthMethodSSH
	case "none", "":
		config.AuthMethod = AuthMethodNone
	default:
		log.Fatalf("Invalid authentication method: %s. Choose from none, https, ssh.", *authMethod)
	}
}

func main() {
	ctx := context.Background()

	if err := promptForMissingInputs(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Welcome to repo-to-txt!")

	repoName, err := extractRepoName(config.RepoURL)
	if err != nil {
		log.Fatalf("Error extracting repository name: %v", err)
	}
	outputFile := fmt.Sprintf("%s%s", repoName, defaultOutputExt)

	tempDir, err := os.MkdirTemp("", defaultCloneDir)
	if err != nil {
		log.Fatalf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	auth, err := setupAuth()
	if err != nil {
		log.Fatalf("Error setting up authentication: %v", err)
	}

	if err := cloneOrPullRepo(ctx, config.RepoURL, tempDir, auth); err != nil {
		log.Fatalf("Error cloning/pulling repository: %v", err)
	}

	if err := writeRepoContentsToFile(tempDir, outputFile, config.ExcludeFolders, config.IncludeExt); err != nil {
		log.Fatalf("Error writing repository contents to file: %v", err)
	}

	log.Printf("Repository contents written to %s", outputFile)
}

// promptForMissingInputs prompts the user for any missing configuration inputs.
func promptForMissingInputs() error {
	if config.RepoURL == "" {
		config.RepoURL = promptForRepoURL()
	}

	if !authFlagSet {
		config.AuthMethod = promptForAuthMethod()
	}

	switch config.AuthMethod {
	case AuthMethodHTTPS:
		if config.Username == "" {
			config.Username = promptForInput("Enter your GitHub username:")
		}
		if config.PersonalAccessToken == "" {
			config.PersonalAccessToken = promptForPassword("Enter your GitHub Personal Access Token:")
		}
	case AuthMethodSSH:
		if config.SSHKeyPath == "" {
			config.SSHKeyPath = promptForInputWithDefault("Enter the path to your SSH private key", defaultSSHKeyPath())
		}
		if isSSHKeyPassphraseProtected(config.SSHKeyPath) {
			config.SSHPassphrase = promptForPassword("Enter your SSH key passphrase (leave empty if none):")
		}
	}

	if len(config.ExcludeFolders) == 0 {
		excludeInput := promptForInputWithDefault("Enter folders to exclude (comma-separated, leave empty to include all)", "")
		config.ExcludeFolders = parseCommaSeparated(excludeInput)
	}

	if len(config.IncludeExt) == 0 {
		includeInput := promptForInputWithDefault("Enter file extensions to include (comma-separated, leave empty to include all)", "")
		config.IncludeExt = parseCommaSeparated(includeInput)
	}

	return validateRepoURL(config.RepoURL, config.AuthMethod)
}

// isSSHKeyPassphraseProtected checks if the SSH key is passphrase protected.
// This is a heuristic and may not be accurate for all keys.
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

// promptForRepoURL prompts the user for a valid repository URL.
func promptForRepoURL() string {
	for {
		input := promptForInput("Enter the GitHub repository URL (HTTPS or SSH):")
		if err := validateRepoURL(input, AuthMethodNone); err == nil {
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
func promptForAuthMethod() AuthMethod {
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
		return AuthMethodHTTPS
	case "SSH":
		return AuthMethodSSH
	default:
		return AuthMethodNone
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

// extractRepoName extracts the repository name from the repository URL.
func extractRepoName(repoURL string) (string, error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("invalid repository URL: %w", err)
	}

	if u.Scheme == "" && strings.Contains(repoURL, "@") && strings.Contains(repoURL, ":") {
		parts := strings.Split(repoURL, ":")
		if len(parts) != 2 {
			return "", errors.New("invalid SSH repository URL format")
		}
		repoPath := strings.TrimSuffix(parts[1], ".git")
		return filepath.Base(repoPath), nil
	}

	repoPath := strings.TrimSuffix(u.Path, ".git")
	repoName := filepath.Base(repoPath)
	if repoName == "" {
		return "", errors.New("could not determine repository name from URL")
	}
	return repoName, nil
}

// validateRepoURL validates the repository URL based on the selected AuthMethod.
func validateRepoURL(repoURL string, _ AuthMethod) error {
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

// setupAuth prepares the authentication method based on the config.
func setupAuth() (transport.AuthMethod, error) {
	switch config.AuthMethod {
	case AuthMethodHTTPS:
		if config.Username == "" || config.PersonalAccessToken == "" {
			return nil, errors.New("username and personal access token must be provided for HTTPS authentication")
		}
		return &http.BasicAuth{
			Username: config.Username, // GitHub username
			Password: config.PersonalAccessToken,
		}, nil
	case AuthMethodSSH:
		if config.SSHPassphrase != "" {
			return ssh.NewPublicKeys(defaultSSHKeyName, []byte(config.SSHPassphrase), config.SSHKeyPath)
		}
		return ssh.NewPublicKeysFromFile(defaultSSHKeyName, config.SSHKeyPath, "")
	case AuthMethodNone:
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported authentication method: %v", config.AuthMethod)
	}
}

// cloneOrPullRepo clones the repository if it doesn't exist locally or pulls the latest changes if it does.
func cloneOrPullRepo(ctx context.Context, repoURL, repoPath string, auth transport.AuthMethod) error {
	log.Printf("Cloning repository: %s", repoURL)
	_, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth:     auth,
	})
	if err != nil {
		// If the repository already exists, attempt to pull the latest changes
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			log.Printf("Repository already exists. Attempting to pull latest changes.")
			repo, err := git.PlainOpen(repoPath)
			if err != nil {
				return fmt.Errorf("failed to open existing repository: %w", err)
			}
			w, err := repo.Worktree()
			if err != nil {
				return fmt.Errorf("failed to get worktree: %w", err)
			}
			err = w.Pull(&git.PullOptions{
				RemoteName: "origin",
				Progress:   os.Stdout,
				Auth:       auth,
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				return fmt.Errorf("failed to pull repository: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

// writeRepoContentsToFile writes the contents of the repository to the specified output file.
func writeRepoContentsToFile(repoPath, outputFile string, excludeFolders, includeExt []string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("unable to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	return filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}

		if shouldExcludeFile(relPath, excludeFolders, includeExt) {
			return nil
		}

		content, err := readFileContent(path)
		if err != nil {
			log.Printf("Skipping file %s: %v", relPath, err)
			return nil
		}

		return writeFileContent(writer, relPath, content)
	})
}

// shouldExcludeFile determines if a file should be excluded based on its path and extension.
func shouldExcludeFile(relPath string, excludeFolders, includeExt []string) bool {
	normalizedRelPath := filepath.ToSlash(relPath)
	for _, exclude := range excludeFolders {
		normalizedExclude := filepath.ToSlash(strings.TrimSpace(exclude))
		if normalizedExclude == "" {
			continue
		}
		if strings.HasPrefix(normalizedRelPath, normalizedExclude+"/") || normalizedRelPath == normalizedExclude {
			return true
		}
	}

	if len(includeExt) > 0 {
		ext := strings.ToLower(filepath.Ext(relPath))
		return !contains(includeExt, ext)
	}

	return strings.HasSuffix(strings.ToLower(relPath), defaultExcludedExt)
}

// readFileContent reads the content of the file if it's a text file.
func readFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if isBinary(buf[:n]) {
		return nil, errors.New("binary file")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

// writeFileContent writes the content of a file to the output writer.
func writeFileContent(writer *bufio.Writer, relPath string, content []byte) error {
	separator := fmt.Sprintf("=== %s ===\n", relPath)
	if _, err := io.WriteString(writer, separator); err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}
	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("error writing file content: %w", err)
	}
	if _, err := io.WriteString(writer, "\n\n"); err != nil {
		return fmt.Errorf("error writing newline to output file: %w", err)
	}
	return nil
}

// isBinary checks if the file content is binary.
func isBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) != -1
}

// contains checks if a slice contains a particular string (case-insensitive).
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// parseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
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
