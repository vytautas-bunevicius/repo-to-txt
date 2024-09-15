// Package output_test contains unit tests for the output package.
package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

// TestWriteRepoContentsToFile verifies that the WriteRepoContentsToFile function
// successfully writes repository contents to an output file based on the configuration.
func TestWriteRepoContentsToFile(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.txt")

	// Create a dummy file in tempDir
	dummyFile := filepath.Join(tempDir, "test.go")
	err := os.WriteFile(dummyFile, []byte("package main\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	cfg := &config.Config{
		ExcludeFolders: []string{},
		IncludeExt:     []string{".go"},
	}

	err = WriteRepoContentsToFile(tempDir, outputFile, cfg)
	if err != nil {
		t.Fatalf("WriteRepoContentsToFile returned an error: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %q does not exist", outputFile)
	}

	// Verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "=== test.go ===\npackage main\n\n\n"
	if string(content) != expectedContent {
		t.Errorf("Output file content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

// TestWriteRepoContentsToFileExclusions verifies that files and folders specified in the
// exclusion list are correctly excluded from the output file.
func TestWriteRepoContentsToFileExclusions(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.txt")

	// Create directories and files
	excludedDir := filepath.Join(tempDir, "docs")
	err := os.Mkdir(excludedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create excluded directory: %v", err)
	}

	excludedFile := filepath.Join(excludedDir, "README.md")
	err = os.WriteFile(excludedFile, []byte("# README"), 0644)
	if err != nil {
		t.Fatalf("Failed to write excluded file: %v", err)
	}

	includedFile := filepath.Join(tempDir, "main.go")
	err = os.WriteFile(includedFile, []byte("package main\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write included file: %v", err)
	}

	cfg := &config.Config{
		ExcludeFolders: []string{"docs"},
		IncludeExt:     []string{".go", ".md"},
	}

	err = WriteRepoContentsToFile(tempDir, outputFile, cfg)
	if err != nil {
		t.Fatalf("WriteRepoContentsToFile returned an error: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %q does not exist", outputFile)
	}

	// Verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "=== main.go ===\npackage main\n\n\n"
	if string(content) != expectedContent {
		t.Errorf("Output file content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

// TestWriteRepoContentsToFileBinaryFile verifies that binary files are skipped.
func TestWriteRepoContentsToFileBinaryFile(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.txt")

	// Create a binary file
	binaryFile := filepath.Join(tempDir, "image.png")
	err := os.WriteFile(binaryFile, []byte{0x89, 0x50, 0x4E, 0x47}, 0644) // PNG header bytes
	if err != nil {
		t.Fatalf("Failed to write binary file: %v", err)
	}

	// Create a text file
	textFile := filepath.Join(tempDir, "main.go")
	err = os.WriteFile(textFile, []byte("package main\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write text file: %v", err)
	}

	cfg := &config.Config{
		ExcludeFolders: []string{},
		IncludeExt:     []string{".go", ".png"},
	}

	err = WriteRepoContentsToFile(tempDir, outputFile, cfg)
	if err != nil {
		t.Fatalf("WriteRepoContentsToFile returned an error: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %q does not exist", outputFile)
	}

	// Verify content contains only the text file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "=== main.go ===\npackage main\n\n\n"
	if string(content) != expectedContent {
		t.Errorf("Output file content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

// TestWriteRepoContentsToFileDefaultExclusions verifies that default excluded extensions are respected.
func TestWriteRepoContentsToFileDefaultExclusions(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.txt")

	// Create a default excluded file
	excludedFile := filepath.Join(tempDir, "notebook.ipynb")
	err := os.WriteFile(excludedFile, []byte("{ \"cells\": [] }"), 0644)
	if err != nil {
		t.Fatalf("Failed to write excluded file: %v", err)
	}

	// Create an included file
	includedFile := filepath.Join(tempDir, "main.go")
	err = os.WriteFile(includedFile, []byte("package main\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write included file: %v", err)
	}

	cfg := &config.Config{
		ExcludeFolders: []string{},
		IncludeExt:     nil, // No specific extensions to include
	}

	err = WriteRepoContentsToFile(tempDir, outputFile, cfg)
	if err != nil {
		t.Fatalf("WriteRepoContentsToFile returned an error: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %q does not exist", outputFile)
	}

	// Verify content contains only the included file
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent := "=== main.go ===\npackage main\n\n\n"
	if string(content) != expectedContent {
		t.Errorf("Output file content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}
