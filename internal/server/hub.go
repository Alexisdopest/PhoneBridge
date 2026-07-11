package server

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message defines the structure for bidirectional websocket sync
type Message struct {
	Version   string      `json:"version"`
	MessageID string      `json:"message_id"`
	Source    string      `json:"source"`
	Event     string      `json:"event"`
	Payload   interface{} `json:"payload"`
	Timestamp int64       `json:"timestamp"`
}

// Client represents a generic connected device (Android/Mac/iOS App)
type Client struct {
	Conn     *websocket.Conn
	Device   string
	SendChan chan []byte
}

// Hub manages active generic Client connections and broadcasts messages
type Hub struct {
	mu      sync.Mutex
	clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
}

// Register adds a new client
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
	log.Printf("WS client registered (Device: %s). Total: %d", client.Device, len(h.clients))
}

// Unregister removes a client safely
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.SendChan)
		client.Conn.Close()
		log.Printf("WS client unregistered (Device: %s). Total: %d", client.Device, len(h.clients))
	}
}

// Broadcast sends a message to all active clients
func (h *Hub) Broadcast(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Broadcast serialization error: %v", err)
		return
	}

	for client := range h.clients {
		select {
		case client.SendChan <- data:
		default:
			// Buffer full, drop client to prevent deadlock
			close(client.SendChan)
			delete(h.clients, client)
			client.Conn.Close()
			log.Printf("WS client dropped due to slow consumer (Device: %s)", client.Device)
		}
	}
}
