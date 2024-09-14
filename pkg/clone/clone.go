package clone

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// ExtractRepoName extracts the repository name from the repository URL.
func ExtractRepoName(repoURL string) (string, error) {
    if strings.HasPrefix(repoURL, "git@") {
        // Handle SSH URLs
        parts := strings.SplitN(repoURL, ":", 2)
        if len(parts) != 2 {
            return "", errors.New("invalid SSH repository URL format")
        }
        path := parts[1]
        repoPath := strings.TrimSuffix(path, ".git")
        repoName := filepath.Base(repoPath)
        if repoName == "" {
            return "", errors.New("could not determine repository name from URL")
        }
        return repoName, nil
    } else if strings.HasPrefix(repoURL, "https://") || strings.HasPrefix(repoURL, "http://") {
        // Handle HTTPS URLs
        u, err := url.Parse(repoURL)
        if err != nil {
            return "", fmt.Errorf("invalid repository URL: %w", err)
        }
        repoPath := strings.TrimSuffix(u.Path, ".git")
        repoName := filepath.Base(repoPath)
        if repoName == "" {
            return "", errors.New("could not determine repository name from URL")
        }
        return repoName, nil
    } else {
        return "", errors.New("invalid repository URL format")
    }
}
// CloneOrPullRepo clones the repository if it doesn't exist locally or pulls the latest changes if it does.
func CloneOrPullRepo(ctx context.Context, repoURL, repoPath string, auth transport.AuthMethod) error {
	fmt.Printf("Cloning repository: %s\n", repoURL)
	_, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth:     auth,
	})
	if err != nil {
		// If the repository already exists, attempt to pull the latest changes
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			fmt.Println("Repository already exists. Attempting to pull latest changes.")
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
