package main

import (
	"log"
	"strconv"

	"github.com/Alexisdopest/PhoneBridge/internal/discovery"
	"github.com/Alexisdopest/PhoneBridge/internal/server"
)

func main() {
	token := "123456" 
	portStr := "8080"
	port, _ := strconv.Atoi(portStr)

	log.Println("Starting PhoneBridge service...")
	log.Printf("Pair Code: %s\n", token)
	log.Println("Please configure your iOS Shortcut to use: Authorization: Bearer 123456")
	
	// Start mDNS
	mdnsServer, err := discovery.StartMDNS(port)
	if err != nil {
		log.Printf("Failed to start mDNS: %v", err)
	} else {
		defer mdnsServer.Shutdown()
	}

	// Start HTTP Server
	if err := server.Start(portStr, token); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
