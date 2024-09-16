// Package output manages the generation of the output text file containing repository contents.
// It handles writing file contents to the output file with appropriate formatting and exclusions,
// as well as copying specified files to the output directory.
package output

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/util"
)

// FindFiles searches for the specified file names within the repository directory.
//
// Parameters:
//   - repoPath: The local path of the cloned repository.
//   - fileNames: A slice of exact file names to search for.
//
// Returns:
//   - map[string][]string: A map where the key is the file name and the value is a slice of matching file paths.
//   - error: An error if the search fails.
func FindFiles(repoPath string, fileNames []string) (map[string][]string, error) {
	if len(fileNames) == 0 {
		return nil, errors.New("no file names provided to search for")
	}

	fileMatches := make(map[string][]string)

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip paths that can't be accessed
			log.Printf("Error accessing path %s: %v", path, err)
			return nil
		}

		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil // Skip directories and hidden files
		}

		for _, fileName := range fileNames {
			if strings.EqualFold(info.Name(), fileName) {
				fileMatches[fileName] = append(fileMatches[fileName], path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %w", repoPath, err)
	}

	return fileMatches, nil
}

// WriteRepoContentsToFile writes the contents of the specified repository directory to an output file.
// It traverses the repository, applies exclusion rules, and formats the output with file separators.
//
// Parameters:
//   - repoPath: The local path of the cloned repository.
//   - outputFile: The path to the output text file.
//   - cfg: A pointer to the Config struct containing exclusion and inclusion rules.
//
// Returns:
//   - error: An error if writing to the file fails.
func WriteRepoContentsToFile(repoPath, outputFile string, cfg *config.Config) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("unable to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil // Skip directories and hidden files
		}

		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}

		if shouldExcludeFile(relPath, cfg) {
			return nil // Skip excluded files
		}

		content, err := readFileContent(path)
		if err != nil {
			log.Printf("Skipping file %s: %v", relPath, err)
			return nil // Skip files that cannot be read or are binary
		}

		return writeFileContent(writer, relPath, content)
	})

	if err != nil {
		return fmt.Errorf("error walking the path %s: %w", repoPath, err)
	}

	return nil
}

// shouldExcludeFile determines whether a file should be excluded based on its relative path and extension.
// It checks against the excluded folders and included extensions specified in the configuration.
//
// Parameters:
//   - relPath: The relative path of the file within the repository.
//   - cfg: A pointer to the Config struct containing exclusion and inclusion rules.
//
// Returns:
//   - bool: True if the file should be excluded, false otherwise.
func shouldExcludeFile(relPath string, cfg *config.Config) bool {
	normalizedRelPath := filepath.ToSlash(relPath)
	for _, exclude := range cfg.ExcludeFolders {
		normalizedExclude := filepath.ToSlash(strings.TrimSpace(exclude))
		if normalizedExclude == "" {
			continue
		}
		if strings.HasPrefix(normalizedRelPath, normalizedExclude+"/") || normalizedRelPath == normalizedExclude {
			return true
		}
	}

	if len(cfg.IncludeExt) > 0 {
		ext := strings.ToLower(filepath.Ext(relPath))
		return !util.Contains(cfg.IncludeExt, ext)
	}

	return strings.HasSuffix(strings.ToLower(relPath), config.DefaultExcludedExt)
}

// readFileContent reads and returns the content of a file if it is a text file.
// It skips binary files by checking for null bytes.
//
// Parameters:
//   - path: The file system path to the file.
//
// Returns:
//   - []byte: The content of the file.
//   - error: An error if the file cannot be read or is identified as binary.
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

// writeFileContent writes the content of a file to the output writer with appropriate formatting.
// It adds a separator with the relative file path before the content.
//
// Parameters:
//   - writer: The buffered writer for the output file.
//   - relPath: The relative path of the file within the repository.
//   - content: The content of the file.
//
// Returns:
//   - error: An error if writing to the output file fails.
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

// isBinary checks if the provided byte slice contains any null bytes, indicating a binary file.
//
// Parameters:
//   - data: The byte slice to check.
//
// Returns:
//   - bool: True if the data is binary, false otherwise.
func isBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) != -1
}
