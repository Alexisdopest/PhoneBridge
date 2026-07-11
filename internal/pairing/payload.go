package pairing

import (
	"encoding/json"
	"net"
)

type QRPayload struct {
	Version string `json:"version"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Device  string `json:"device"`
	Token   string `json:"token"`
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func GeneratePayload(port, token string) string {
	payload := QRPayload{
		Version: "1.1",
		Host:    GetLocalIP(),
		Port:    port,
		Device:  "windows_host",
		Token:   token,
	}
	data, _ := json.Marshal(payload)
	return string(data)
}
