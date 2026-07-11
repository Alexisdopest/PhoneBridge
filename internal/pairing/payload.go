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
	interfaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}

	var fallbackIP string

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					ipStr := ip4.String()
					// 过滤掉 APIPA 无效网段 (如 Tailscale 未分配、本地连接回环)
					if ip4[0] == 169 && ip4[1] == 254 {
						continue
					}
					// 优先匹配局域网私有网段
					if (ip4[0] == 192 && ip4[1] == 168) ||
						(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
						(ip4[0] == 10) {
						return ipStr
					}
					fallbackIP = ipStr
				}
			}
		}
	}

	if fallbackIP != "" {
		return fallbackIP
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
