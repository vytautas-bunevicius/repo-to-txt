// Package main serves as the entry point for the repo-to-txt CLI tool.
// It orchestrates the configuration parsing, user prompting, repository cloning/pulling,
// and the generation of a text file containing the repository's contents.
package main

import (
	"bufio" // Added bufio package
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
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
//  9. Writes the repository contents or copies specified files to the output directory.
//
// 10. Optionally copies the contents to clipboard if requested.
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
		return fmt.Errorf("error prompting for inputs: %w", err)
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

	if len(cfg.FileNames) > 0 {
		// Handle writing specified files' contents to outputFile
		fileMatches, err := output.FindFiles(tempDir, cfg.FileNames)
		if err != nil {
			return fmt.Errorf("error searching for specified files: %w", err)
		}

		// Open outputFile for writing (truncate if exists)
		outFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("unable to create output file: %w", err)
		}
		defer outFile.Close()

		writer := bufio.NewWriter(outFile)
		defer writer.Flush()

		// Iterate over each specified file name
		for _, fileName := range cfg.FileNames {
			matches, exists := fileMatches[fileName]
			if !exists || len(matches) == 0 {
				log.Printf("No matches found for file name: %s", fileName)
				continue
			}

			var selectedPath string
			if len(matches) == 1 {
				selectedPath = matches[0]
				log.Printf("Found one match for %s: %s", fileName, selectedPath)
			} else {
				// Prompt user to select which file to include
				selectedPath, err = prompt.SelectFile(fileName, matches)
				if err != nil {
					return fmt.Errorf("error selecting file for %s: %w", fileName, err)
				}
			}

			// Read file content
			content, err := ioutil.ReadFile(selectedPath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", selectedPath, err)
				continue
			}

			// Compute relative path
			relPath, err := filepath.Rel(tempDir, selectedPath)
			if err != nil {
				log.Printf("Failed to compute relative path for %s: %v", selectedPath, err)
				relPath = filepath.Base(selectedPath) // fallback to base name
			}

			// Write content to outputFile with separator
			separator := fmt.Sprintf("=== %s ===\n", relPath) // include relative path
			if _, err := writer.WriteString(separator); err != nil {
				log.Printf("Failed to write separator to output file: %v", err)
				continue
			}
			if _, err := writer.Write(content); err != nil {
				log.Printf("Failed to write content to output file: %v", err)
				continue
			}
			if _, err := writer.WriteString("\n\n"); err != nil {
				log.Printf("Failed to write newline to output file: %v", err)
				continue
			}

			log.Printf("Added %s to %s", selectedPath, outputFile)
		}

		// After writing all specified files
		log.Printf("Specified files' contents written to %s", outputFile)

		// Handle clipboard copy if requested
		if cfg.CopyToClipboard {
			content, err := ioutil.ReadFile(outputFile)
			if err != nil {
				return fmt.Errorf("error reading output file for clipboard: %w", err)
			}
			if err := clipboard.WriteAll(string(content)); err != nil {
				return fmt.Errorf("error copying to clipboard: %w", err)
			}
			log.Println("Specified files' contents have been copied to the clipboard.")
		} else {
			log.Println("Specified files' contents were not copied to the clipboard.")
		}
	} else {
		// Write the repository contents to the specified output file.
		if err := output.WriteRepoContentsToFile(tempDir, outputFile, cfg); err != nil {
			return fmt.Errorf("error writing repository contents to file: %w", err)
		}
		log.Printf("Repository contents written to %s", outputFile)

		// Handle clipboard copy if requested
		if cfg.CopyToClipboard {
			content, err := ioutil.ReadFile(outputFile)
			if err != nil {
				return fmt.Errorf("error reading output file for clipboard: %w", err)
			}
			if err := clipboard.WriteAll(string(content)); err != nil {
				return fmt.Errorf("error copying to clipboard: %w", err)
			}
			log.Println("Repository contents have been copied to the clipboard.")
		} else {
			log.Println("Repository contents were not copied to the clipboard.")
		}
	}

	return nil
}
