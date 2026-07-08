package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Alexisdopest/PhoneBridge/internal/clipboard"
	"github.com/Alexisdopest/PhoneBridge/internal/storage"
)

// ClipboardHandler handles POST requests to update the clipboard
func ClipboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the text from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	text := string(body)
	if text == "" {
		http.Error(w, "Empty content", http.StatusBadRequest)
		return
	}

	// Write to Windows clipboard
	err = clipboard.WriteText(text)
	if err != nil {
		log.Printf("Failed to write clipboard: %v", err)
		http.Error(w, "Failed to write clipboard", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully wrote %d bytes to clipboard", len(text))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Clipboard updated"))
}

// UploadHandler handles POST requests for uploading files/images
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit to 100 MB max for now to avoid memory explosion
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form (max 100MB)", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing 'file' field in form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine save directory (Downloads/PhoneBridge for now)
	homeDir, _ := os.UserHomeDir()
	destDir := filepath.Join(homeDir, "Downloads", "PhoneBridge")

	savedPath, err := storage.SaveFile(destDir, header.Filename, file)
	if err != nil {
		log.Printf("Failed to save file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully saved file: %s", savedPath)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully: " + savedPath))
}
