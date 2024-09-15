// Package main_test contains integration tests for the repo-to-txt CLI tool.
// These tests verify the end-to-end functionality of the tool, including building the binary,
// cloning a repository, generating the output text file, and optionally copying to the clipboard.
//
// Note: Integration tests are skipped by default and can be enabled by setting the
// RUN_INTEGRATION_TESTS environment variable to "true". Clipboard tests require
// the RUN_CLIPBOARD_TESTS environment variable to be set to "true" and may fail
// in headless or restricted environments.
package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/atotto/clipboard"
)

// TestMainIntegration performs an end-to-end integration test of the repo-to-txt CLI tool.
// It builds the binary, clones a specified repository without authentication, generates the output
// text file, and optionally verifies clipboard copying.
//
// Preconditions:
//   - The RUN_INTEGRATION_TESTS environment variable must be set to "true" to run this test.
//   - Network access is required to clone the specified GitHub repository.
//   - Clipboard access is required to verify clipboard copying if RUN_CLIPBOARD_TESTS is set to "true".
//
// Steps:
//  1. Skip the test if running in short mode or if RUN_INTEGRATION_TESTS is not set to "true".
//  2. Build the repo-to-txt binary.
//  3. Create a temporary output directory.
//  4. Execute the binary with specified command-line arguments, including the clipboard flag.
//  5. Verify that the output file is created and contains content.
//  6. Optionally, verify that the clipboard contains the expected content.
//
// Parameters:
//   - t: The testing framework's testing object.
func TestMainIntegration(t *testing.T) {
	// Skip the test in short mode to allow faster test runs.
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check for the RUN_INTEGRATION_TESTS environment variable to determine if integration tests should run.
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	// Step 1: Build the binary.
	exePath := filepath.Join(os.TempDir(), "repo-to-txt-test")
	buildCmd := exec.Command("go", "build", "-o", exePath, "./")
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, string(buildOutput))
	}
	defer os.Remove(exePath) // Clean up the binary after the test.

	// Step 2: Create a temporary output directory.
	outputDir, err := ioutil.TempDir("", "repo-to-txt-output")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir) // Clean up the output directory after the test.

	// Step 3: Prepare command-line arguments.
	repoURL := "https://github.com/vytautas-bunevicius/repo-to-txt.git"
	args := []string{
		"-repo=" + repoURL,
		"-auth=none",
		"-output-dir=" + outputDir,
		"-copy-clipboard=true", // Enable clipboard copying
	}

	// Step 4: Run the binary.
	runCmd := exec.Command(exePath, args...)

	// Prepare input for potential prompts (e.g., folders to exclude).
	input := bytes.NewBufferString("\n") // Provide an empty input for excluded folders.
	runCmd.Stdin = input

	// Capture standard output and standard error.
	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	// Execute the command.
	if err := runCmd.Run(); err != nil {
		t.Fatalf("Command failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	// Step 5: Verify that the output file exists.
	outputFile := filepath.Join(outputDir, "repo-to-txt.txt")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to exist", outputFile)
	}

	// Step 6: Optionally, read and verify the contents of the output file.
	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Output file %s is empty", outputFile)
	}

	// New Step: Verify clipboard content (optional and environment-dependent)
	if os.Getenv("RUN_CLIPBOARD_TESTS") == "true" {
		clipboardContent, err := clipboard.ReadAll()
		if err != nil {
			t.Fatalf("Failed to read from clipboard: %v", err)
		}

		// Simple comparison; in real tests, consider more robust checks
		if !strings.Contains(string(clipboardContent), "=== repo-to-txt.txt ===") {
			t.Errorf("Clipboard content does not contain expected output")
		}
	}
}
