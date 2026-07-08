package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Alexisdopest/PhoneBridge/internal/clipboard"
	"github.com/Alexisdopest/PhoneBridge/internal/storage"
	"github.com/gorilla/websocket"
)

// ClipboardHandler handles POST requests to update the clipboard
func ClipboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for LAN connections initially
	},
}

// WSEventsHandler upgrades the HTTP connection to a WebSocket for future event pushing
func WSEventsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket client connected (Pre-reserved for Milestone 3+)")
	
	// Keep connection alive until client disconnects
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	log.Println("WebSocket client disconnected")
}
