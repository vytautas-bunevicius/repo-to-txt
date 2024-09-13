package main

import (
	"bufio"
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
)

const (
	defaultCloneDir    = "repo-to-txt-clone"
	defaultOutputExt   = ".txt"
	defaultSSHKeyName  = "git"
	defaultExcludedExt = ".ipynb" // Default exclusion for Jupyter notebooks
)

// AuthMethod represents the type of authentication to use.
type AuthMethod int

const (
	AuthMethodNone AuthMethod = iota
	AuthMethodHTTPS
	AuthMethodSSH
)

// String returns the string representation of the AuthMethod.
func (a AuthMethod) String() string {
	return [...]string{"No Authentication", "HTTPS with PAT", "SSH"}[a]
}

// Config holds the configuration for the CLI tool.
type Config struct {
	RepoURL             string
	AuthMethod          AuthMethod
	Username            string
	PersonalAccessToken string
	SSHKeyPath          string
	ExcludeFolders      []string
	IncludeExt          []string
}

// parseFlags parses command-line flags and returns a Config struct.
func parseFlags() Config {
	var config Config
	var excludeFolders string
	var includeExt string

	flag.StringVar(&config.RepoURL, "repo", "", "GitHub repository URL (HTTPS or SSH) (Required)")
	authMethod := flag.String("auth", "", "Authentication method: none, https, or ssh (Required)")
	flag.StringVar(&config.Username, "username", "", "GitHub username (for HTTPS)")
	flag.StringVar(&config.PersonalAccessToken, "pat", "", "GitHub Personal Access Token (for HTTPS)")
	flag.StringVar(&config.SSHKeyPath, "ssh-key", "", "Path to SSH private key (for SSH)")
	flag.StringVar(&excludeFolders, "exclude", "", "Comma-separated list of folders to exclude from the output")
	flag.StringVar(&includeExt, "include-ext", "", "Comma-separated list of file extensions to include (e.g., .go,.md). If not set, defaults to excluding certain non-code files like .ipynb")

	flag.Parse()

	if excludeFolders != "" {
		config.ExcludeFolders = parseCommaSeparated(excludeFolders)
	}

	if includeExt != "" {
		config.IncludeExt = parseCommaSeparated(includeExt)
	} else {
		// Default exclusion for non-code files
		config.IncludeExt = []string{}
	}

	// Set AuthMethod based on flag
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

	return config
}

// parseCommaSeparated splits a comma-separated string into a slice of trimmed strings.
func parseCommaSeparated(input string) []string {
	parts := strings.Split(input, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// promptForMissingInputs prompts the user for any missing configuration inputs.
func promptForMissingInputs(config *Config) error {
	// Prompt for RepoURL if not provided
	if config.RepoURL == "" {
		config.RepoURL = promptForInput("Enter the GitHub repository URL (HTTPS or SSH):")
	}

	// Prompt for AuthMethod if not set via flags
	if config.AuthMethod == AuthMethodNone && config.AuthMethod != AuthMethodHTTPS && config.AuthMethod != AuthMethodSSH {
		config.AuthMethod = promptForAuthMethod()
	}

	switch config.AuthMethod {
	case AuthMethodHTTPS:
		// Prompt for Username if not provided
		if config.Username == "" {
			config.Username = promptForInput("Enter your GitHub username:")
		}

		// Prompt for Personal Access Token if not provided
		if config.PersonalAccessToken == "" {
			config.PersonalAccessToken = promptForPassword("Enter your GitHub Personal Access Token:")
		}

	case AuthMethodSSH:
		// Prompt for SSHKeyPath if not provided
		if config.SSHKeyPath == "" {
			config.SSHKeyPath = promptForInputWithDefault("Enter the path to your SSH private key", defaultSSHKeyPath())
		}

	case AuthMethodNone:
		// No additional input required
	default:
		return fmt.Errorf("unsupported authentication method: %s", config.AuthMethod)
	}

	// Prompt for ExcludeFolders if not provided via flags
	if len(config.ExcludeFolders) == 0 {
		excludeInput := promptForInputWithDefault("Enter folders to exclude (comma-separated, leave empty to include all)", "")
		if excludeInput != "" {
			config.ExcludeFolders = parseCommaSeparated(excludeInput)
		}
	}

	return nil
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
	prompt := promptui.Prompt{
		Label:    label,
		Mask:     '*',
		Validate: validateNonEmpty,
		Stdin:    os.Stdin,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Password prompt failed: %v", err)
	}
	return result
}

// validateNonEmpty ensures the input is not empty.
func validateNonEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
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

	// Handle SSH URLs like git@github.com:user/repo.git
	if u.Scheme == "" && strings.Contains(repoURL, "@") && strings.Contains(repoURL, ":") {
		parts := strings.Split(repoURL, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid SSH repository URL format")
		}
		repoPath := parts[1]
		repoPath = strings.TrimSuffix(repoPath, ".git")
		repoName := filepath.Base(repoPath)
		return repoName, nil
	}

	// Handle standard URLs
	repoPath := strings.TrimSuffix(u.Path, ".git")
	repoName := filepath.Base(repoPath)
	if repoName == "" {
		return "", fmt.Errorf("could not determine repository name from URL")
	}
	return repoName, nil
}

// cloneOrPullRepoHTTPS handles cloning or pulling the repository using HTTPS with PAT.
func cloneOrPullRepoHTTPS(repoURL, repoPath, username, pat string) error {
	auth := &http.BasicAuth{
		Username: username, // can be anything except an empty string
		Password: pat,
	}

	return cloneOrPullRepo(repoURL, repoPath, auth)
}

// cloneOrPullRepoSSH handles cloning or pulling the repository using SSH authentication.
func cloneOrPullRepoSSH(repoURL, repoPath, sshKeyPath string) error {
	auth, err := ssh.NewPublicKeysFromFile(defaultSSHKeyName, sshKeyPath, "")
	if err != nil {
		return fmt.Errorf("error creating SSH auth method: %w", err)
	}

	return cloneOrPullRepo(repoURL, repoPath, auth)
}

// cloneOrPullRepo clones the repository if it doesn't exist locally or pulls the latest changes if it does.
func cloneOrPullRepo(repoURL, repoPath string, auth transport.AuthMethod) error {
	_, err := os.Stat(repoPath)
	if os.IsNotExist(err) {
		log.Printf("Cloning repository %s...\n", repoURL)
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
			Auth:     auth,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	} else if err == nil {
		log.Printf("Pulling latest changes from %s...\n", repoURL)
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}
		worktree, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %w", err)
		}
		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			Auth:       auth,
		})
		if err != nil {
			if errors.Is(err, git.NoErrAlreadyUpToDate) {
				log.Println("Repository is already up to date")
				return nil
			}
			return fmt.Errorf("failed to pull repository: %w", err)
		}
	} else {
		return fmt.Errorf("error checking repository path: %w", err)
	}
	return nil
}

// writeRepoContentsToFile writes the contents of the repository to the specified output file.
func writeRepoContentsToFile(repoPath, outputFile string, excludeFolders []string, includeExt []string) error {
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
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}

		// Check if the file is in any of the excluded folders
		for _, exclude := range excludeFolders {
			exclude = strings.TrimSpace(exclude)
			if exclude == "" {
				continue
			}
			// Normalize paths to use forward slashes
			normalizedRelPath := filepath.ToSlash(relPath)
			normalizedExclude := filepath.ToSlash(exclude)
			if strings.HasPrefix(normalizedRelPath, normalizedExclude+"/") || normalizedRelPath == normalizedExclude {
				// Skip this file
				return nil
			}
		}

		// Handle file extensions
		if len(includeExt) > 0 {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if !contains(includeExt, ext) {
				// Skip files that do not match the included extensions
				return nil
			}
		} else {
			// Default exclusion for non-code files
			if strings.HasSuffix(strings.ToLower(info.Name()), defaultExcludedExt) {
				return nil
			}
		}

		content, err := readFileContent(path)
		if err != nil {
			// Skip files that can't be read (e.g., binary files)
			log.Printf("Skipping file %s: %v\n", relPath, err)
			return nil
		}

		separator := fmt.Sprintf("=== %s ===\n", relPath)
		if _, err := writer.WriteString(separator); err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
		if _, err := writer.Write(content); err != nil {
			return fmt.Errorf("error writing file content: %w", err)
		}
		if _, err := writer.WriteString("\n\n"); err != nil {
			return fmt.Errorf("error writing newline to output file: %w", err)
		}

		return nil
	})
}

// readFileContent reads the content of the file if it's a text file.
func readFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Simple check to skip binary files
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if isBinary(buf[:n]) {
		return nil, fmt.Errorf("binary file")
	}

	// Reset the file pointer to the beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

// isBinary checks if the file content is binary.
func isBinary(data []byte) bool {
	for _, b := range data {
		// If the byte is a null byte, it's likely binary
		if b == 0 {
			return true
		}
	}
	return false
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

func main() {
	// Parse command-line flags
	config := parseFlags()

	// Validate required flags
	if config.RepoURL == "" {
		log.Fatal("Repository URL must be provided using the -repo flag or interactively.")
	}

	// Prompt for missing configuration inputs interactively
	if err := promptForMissingInputs(&config); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Welcome to repo-to-txt!")

	// Extract repository name to name the output file
	repoName, err := extractRepoName(config.RepoURL)
	if err != nil {
		log.Fatalf("Error extracting repository name: %v", err)
	}
	outputFile := fmt.Sprintf("%s%s", repoName, defaultOutputExt)

	// Use a unique temporary directory to avoid conflicts
	tempDir, err := os.MkdirTemp("", defaultCloneDir)
	if err != nil {
		log.Fatalf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up the cloned repository after execution

	// Handle cloning or pulling based on authentication method
	switch config.AuthMethod {
	case AuthMethodHTTPS:
		err = cloneOrPullRepoHTTPS(config.RepoURL, tempDir, config.Username, config.PersonalAccessToken)
	case AuthMethodSSH:
		err = cloneOrPullRepoSSH(config.RepoURL, tempDir, config.SSHKeyPath)
	case AuthMethodNone:
		err = cloneOrPullRepo(config.RepoURL, tempDir, nil)
	default:
		log.Fatalf("Unsupported authentication method: %s", config.AuthMethod)
	}

	if err != nil {
		log.Fatalf("Error cloning/pulling repository: %v", err)
	}

	// Write repository contents to the output file, excluding specified folders
	err = writeRepoContentsToFile(tempDir, outputFile, config.ExcludeFolders, config.IncludeExt)
	if err != nil {
		log.Fatalf("Error writing repository contents to file: %v", err)
	}

	log.Printf("Repository contents written to %s\n", outputFile)
}
