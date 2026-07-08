package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// SaveFile securely saves the file stream to the specified directory, avoiding name collisions.
func SaveFile(destDir string, filename string, r io.Reader) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	if base == "" {
		base = "upload"
	}
	
	// Create a unique filename with timestamp to prevent collisions
	finalName := fmt.Sprintf("%s_%d%s", base, time.Now().UnixMilli(), ext)
	outPath := filepath.Join(destDir, finalName)

	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, r); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return outPath, nil
}
