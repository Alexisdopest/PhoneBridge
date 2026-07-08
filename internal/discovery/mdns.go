package discovery

import (
	"log"
	"os"

	"github.com/hashicorp/mdns"
)

// StartMDNS starts the mDNS server to broadcast the PhoneBridge service.
func StartMDNS(port int) (*mdns.Server, error) {
	host, _ := os.Hostname()
	info := []string{"PhoneBridge Service"}
	
	// Broadcast as "_phonebridge._tcp"
	service, err := mdns.NewMDNSService(host, "_phonebridge._tcp", "", "", port, nil, info)
	if err != nil {
		return nil, err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}
	
	log.Printf("mDNS discovery started for service _phonebridge._tcp on port %d", port)
	return server, nil
}
