package server

import (
	"io"
	"log"
	"net/http"

	"github.com/Alexisdopest/PhoneBridge/internal/clipboard"
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
