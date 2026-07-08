package app

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
	"github.com/Alexisdopest/PhoneBridge/internal/tray"
	"github.com/hashicorp/mdns"
)

type App struct {
	portStr    string
	token      string
	httpServer *http.Server
	mdnsServer *mdns.Server
	trayMgr    *tray.TrayManager
	quitChan   chan struct{}
}

func NewApp() *App {
	return &App{
		portStr:  "8080",
		token:    "123456",
		quitChan: make(chan struct{}),
	}
}

func (a *App) Run() {
	port, _ := strconv.Atoi(a.portStr)

	log.Println("Starting PhoneBridge service...")
	log.Printf("Pair Code: %s\n", a.token)

	// 1. Start HTTP Server
	a.httpServer = server.NewServer(a.portStr, a.token)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server error: %v", err)
		}
	}()

	// 2. Start mDNS
	var err error
	a.mdnsServer, err = discovery.StartMDNS(port)
	if err != nil {
		log.Printf("Failed to start mDNS: %v", err)
	}

	// 3. Setup OS Signal trapping
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 4. Start Tray Manager (Must block main thread)
	a.trayMgr = tray.NewTrayManager(func() {
		// Stop the tray loop and close quit channel
		a.trayMgr.Stop()
		close(a.quitChan)
	})

	go func() {
		select {
		case <-sigChan:
			log.Println("Received OS interrupt signal")
			a.trayMgr.Stop() // This will break the tray.Start() blocking call
			// Only close quitChan if it hasn't been closed
			select {
			case <-a.quitChan:
			default:
				close(a.quitChan)
			}
		case <-a.quitChan:
			log.Println("Tray requested application exit")
		}
	}()

	// This blocks the main thread
	a.trayMgr.Start()

	// 5. Graceful Shutdown phase starts here (after tray exits)
	a.shutdown()
}

func (a *App) shutdown() {
	log.Println("Shutting down servers gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if a.mdnsServer != nil {
		a.mdnsServer.Shutdown()
		log.Println("mDNS server stopped")
	}

	if a.httpServer != nil {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP Server forced to shutdown: %v", err)
		}
	}

	log.Println("PhoneBridge exited")
}
