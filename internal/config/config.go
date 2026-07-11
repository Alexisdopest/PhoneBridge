package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Token string `json:"token"`
	Port  string `json:"port"`
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func LoadConfig() *Config {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".phonebridge")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.json")

	var cfg Config
	data, err := os.ReadFile(configPath)
	if err == nil {
		// Verify token is robust (at least 64 hex chars = 32 bytes)
		if err := json.Unmarshal(data, &cfg); err == nil && len(cfg.Token) >= 64 {
			return &cfg
		}
	}

	cfg = Config{
		Token: generateToken(),
		Port:  "8080",
	}

	data, _ = json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(configPath, data, 0644)
	return &cfg
}
