package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vytautas-bunevicius/repo-to-txt/pkg/auth"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/clone"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/output"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/prompt"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(ctx context.Context) error {
	cfg := config.NewConfig()

	if err := cfg.ParseFlags(); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	if err := prompt.PromptForMissingInputs(cfg); err != nil {
		return err
	}

	log.Println("Welcome to repo-to-txt!")

	repoName, err := clone.ExtractRepoName(cfg.RepoURL)
	if err != nil {
		return fmt.Errorf("error extracting repository name: %w", err)
	}

	outputFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s%s", repoName, config.DefaultOutputExt))

	tempDir, err := os.MkdirTemp("", config.DefaultCloneDir)
	if err != nil {
		return fmt.Errorf("unable to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	authMethod, err := auth.SetupAuth(cfg)
	if err != nil {
		return fmt.Errorf("error setting up authentication: %w", err)
	}

	if err := clone.CloneOrPullRepo(ctx, cfg.RepoURL, tempDir, authMethod); err != nil {
		return fmt.Errorf("error cloning/pulling repository: %w", err)
	}

	if err := output.WriteRepoContentsToFile(tempDir, outputFile, cfg); err != nil {
		return fmt.Errorf("error writing repository contents to file: %w", err)
	}

	log.Printf("Repository contents written to %s", outputFile)
	return nil
}
