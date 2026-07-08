package storage

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSaveFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
		content  string
	}{
		{"Normal file", "test.txt", "hello world"},
		{"Path traversal", "../../../windows/system32/hack.exe", "malicious"},
		{"Absolute path", "C:\\hack.bat", "bad"},
		{"No extension", "testfile", "no ext"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(tt.content))
			savedPath, err := SaveFile(tmpDir, tt.filename, r)
			if err != nil {
				t.Fatalf("SaveFile failed: %v", err)
			}

			// Verify it's securely within tmpDir (no traversal)
			if !strings.HasPrefix(savedPath, tmpDir) {
				t.Errorf("Path traversal vulnerability: file saved outside target dir: %s", savedPath)
			}

			// Verify content
			data, err := os.ReadFile(savedPath)
			if err != nil {
				t.Fatalf("Failed to read saved file: %v", err)
			}
			if string(data) != tt.content {
				t.Errorf("Expected content %s, got %s", tt.content, string(data))
			}
		})
	}
}
