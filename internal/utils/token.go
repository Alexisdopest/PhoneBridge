package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken generates a random hex string of given length bytes.
func GenerateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "123456" // Fallback
	}
	return hex.EncodeToString(b)
}
