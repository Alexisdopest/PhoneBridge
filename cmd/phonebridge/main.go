package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
	}

	// Setup HTTP Server
	srv := server.NewServer(portStr, token)

	// Start server in a goroutine so it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers gracefully...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if mdnsServer != nil {
		mdnsServer.Shutdown()
		log.Println("mDNS server stopped")
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP Server forced to shutdown: %v", err)
	}

	log.Println("PhoneBridge exited")
}
