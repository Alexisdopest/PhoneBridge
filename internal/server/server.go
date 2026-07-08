package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alexisdopest/PhoneBridge/internal/auth"
)

// NewServer initializes routing and returns an HTTP server instance
func NewServer(port string, token string) *http.Server {
	mux := http.NewServeMux()

	// Register routes with auth middleware
	mux.HandleFunc("/api/clipboard", auth.Middleware(token, ClipboardHandler))
	mux.HandleFunc("/api/upload", auth.Middleware(token, UploadHandler))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on http://0.0.0.0%s", addr)
	
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
