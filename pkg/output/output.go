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

// WriteRepoContentsToFile writes the contents of the repository to the specified output file.
func WriteRepoContentsToFile(repoPath, outputFile string, cfg *config.Config) error {
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
}

// shouldExcludeFile determines if a file should be excluded based on its path and extension.
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

// readFileContent reads the content of the file if it's a text file.
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

// writeFileContent writes the content of a file to the output writer.
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

// isBinary checks if the file content is binary.
func isBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) != -1
}
