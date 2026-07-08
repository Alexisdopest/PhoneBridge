package main

import (
	"log"

	"github.com/Alexisdopest/PhoneBridge/internal/server"
)

func main() {
	// For Milestone 1, we use a fixed static token for easy iOS shortcut testing.
	// In the future (v1.1), we can use utils.GenerateToken(16) and display it on the UI.
	token := "123456" 
	port := "8080"

	log.Println("Starting PhoneBridge service...")
	log.Printf("Pair Code: %s\n", token)
	log.Println("Please configure your iOS Shortcut to use: Authorization: Bearer 123456")
	
	if err := server.Start(port, token); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
