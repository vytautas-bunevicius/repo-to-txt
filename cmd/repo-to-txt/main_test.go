package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check for a specific environment variable to run this test
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	// Rest of the test remains the same...
	// Build the binary
	exePath := filepath.Join(os.TempDir(), "repo-to-txt-test")
	cmd := exec.Command("go", "build", "-o", exePath, "./")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, string(output))
	}
	defer os.Remove(exePath)

	// Set up a temporary output directory
	outputDir, err := ioutil.TempDir("", "repo-to-txt-output")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Prepare command-line arguments
	repoURL := "https://github.com/vytautas-bunevicius/repo-to-txt.git"
	args := []string{
		"-repo=" + repoURL,
		"-auth=none",
		"-output-dir=" + outputDir,
	}

	// Run the binary
	cmd = exec.Command(exePath, args...)

	// Prepare input for potential prompts
	input := bytes.NewBufferString("\n") // Empty input for folders to exclude
	cmd.Stdin = input

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		t.Fatalf("Command failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	// Verify the output file exists
	outputFile := filepath.Join(outputDir, "repo-to-txt.txt")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to exist", outputFile)
	}

	// Optionally, read and verify the contents of the output file
	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Output file %s is empty", outputFile)
	}
}
