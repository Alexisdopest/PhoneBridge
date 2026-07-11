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

// UploadHandler handles file and image uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// [Security P1 Fix]: Limit upload size to 1GB to prevent LAN abuse / disk exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 1<<30)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form or file too large", http.StatusBadRequest)
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

// WSEventsHandler upgrades the HTTP connection to a WebSocket for bidirectional event pushing
func WSEventsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	// Extract device from query params if available
	device := r.URL.Query().Get("device")
	if device == "" {
		device = "unknown_client"
	}

	client := &Client{
		Conn:     conn,
		Device:   device,
		SendChan: make(chan []byte, 256),
	}

	hub.Register(client)

	// writePump: handles outgoing messages to this client
	go func() {
		defer client.Conn.Close()
		for message := range client.SendChan {
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}()

	// readPump: handles incoming WS messages
	defer hub.Unregister(client)
	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		
		log.Printf("Received WS message [ID: %s] from %s: %s", msg.MessageID, msg.Source, msg.Event)
		if msg.Event == "clipboard_sync" {
			if payloadMap, ok := msg.Payload.(map[string]interface{}); ok {
				if payloadMap["type"] == "text" {
					if text, ok := payloadMap["data"].(string); ok {
						clipboard.WriteText(text)
						log.Printf("Synchronized clipboard text from WS %s", msg.Source)
					}
				}
			}
		}
	}
}
