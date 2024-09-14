package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

func TestWriteRepoContentsToFile(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.txt")

	// Create a dummy file in tempDir
	dummyFile := filepath.Join(tempDir, "test.go")
	os.WriteFile(dummyFile, []byte("package main\n"), 0644)

	cfg := &config.Config{
		ExcludeFolders: []string{},
		IncludeExt:     []string{".go"},
	}

	err := WriteRepoContentsToFile(tempDir, outputFile, cfg)
	if err != nil {
		t.Fatalf("WriteRepoContentsToFile returned an error: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %q does not exist", outputFile)
	}
}
