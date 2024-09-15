// Package main serves as the entry point for the repo-to-txt CLI tool.
// It orchestrates the configuration parsing, user prompting, repository cloning/pulling,
// the generation of a text file containing the repository's contents, and optionally
// copies the output to the clipboard.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard" // Import for clipboard operations
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/auth"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/clone"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/output"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/prompt"
)

// main is the entry point of the repo-to-txt application.
// It initializes the context and invokes the run function.
// Any errors encountered during execution are logged and cause the program to exit.
func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// run orchestrates the main workflow of the repo-to-txt tool.
// It performs the following steps:
//  1. Initializes a new configuration instance.
//  2. Parses command-line flags into the configuration.
//  3. Prompts the user for any missing configuration inputs.
//  4. Extracts the repository name from the provided URL.
//  5. Determines the output file path based on the configuration.
//  6. Creates a temporary directory for cloning the repository.
//  7. Sets up the authentication method.
//  8. Clones the repository or pulls the latest changes if it already exists locally.
//  9. Writes the repository contents to the specified output file.
//
// 10. Optionally copies the output to the clipboard based on the configuration.
//
// Parameters:
//   - ctx: The context for managing cancellation and deadlines.
//
// Returns:
//   - error: An error if any step in the workflow fails.
func run(ctx context.Context) error {
	// Initialize a new configuration instance.
	cfg := config.NewConfig()

	// Parse command-line flags into the configuration.
	if err := cfg.ParseFlags(); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	// Prompt the user for any missing configuration inputs.
	if err := prompt.PromptForMissingInputs(cfg); err != nil {
		return err
	}

	log.Println("Welcome to repo-to-txt!")

	// Extract the repository name from the provided URL.
	repoName, err := clone.ExtractRepoName(cfg.RepoURL)
	if err != nil {
		return fmt.Errorf("error extracting repository name: %w", err)
	}

	// Determine the output file path based on the configuration.
	outputFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s%s", repoName, config.DefaultOutputExt))

	// Create a temporary directory for cloning the repository.
	tempDir, err := os.MkdirTemp("", config.DefaultCloneDir)
	if err != nil {
		return fmt.Errorf("unable to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Ensure the temporary directory is removed after execution.

	// Set up the authentication method based on the configuration.
	authMethod, err := auth.SetupAuth(cfg)
	if err != nil {
		return fmt.Errorf("error setting up authentication: %w", err)
	}

	// Clone the repository or pull the latest changes if it already exists locally.
	if err := clone.CloneOrPullRepo(ctx, cfg.RepoURL, tempDir, authMethod); err != nil {
		return fmt.Errorf("error cloning/pulling repository: %w", err)
	}

	// Write the repository contents to the specified output file.
	if err := output.WriteRepoContentsToFile(tempDir, outputFile, cfg); err != nil {
		return fmt.Errorf("error writing repository contents to file: %w", err)
	}

	log.Printf("Repository contents written to %s", outputFile)

	// Copy to clipboard if requested.
	if cfg.CopyToClipboard {
		content, err := os.ReadFile(outputFile)
		if err != nil {
			return fmt.Errorf("failed to read output file for clipboard copy: %w", err)
		}

		// Attempt to copy to clipboard.
		err = clipboard.WriteAll(string(content))
		if err != nil {
			return fmt.Errorf("failed to copy content to clipboard: %w", err)
		}

		log.Println("Repository contents have been copied to the clipboard.")
	}

	return nil
}
