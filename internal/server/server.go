package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alexisdopest/PhoneBridge/internal/auth"
)

// Start initializes routing and starts the HTTP server
func Start(port string, token string) error {
	mux := http.NewServeMux()

	// Register routes with auth middleware
	mux.HandleFunc("/api/clipboard", auth.Middleware(token, ClipboardHandler))
	mux.HandleFunc("/api/upload", auth.Middleware(token, UploadHandler))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on http://0.0.0.0%s", addr)
	
	// Start server
	return http.ListenAndServe(addr, mux)
}
