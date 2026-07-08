package utils

import "testing"

func TestGenerateToken(t *testing.T) {
	token1 := GenerateToken(16)
	token2 := GenerateToken(16)

	if token1 == "" || token2 == "" {
		t.Error("Generated token should not be empty")
	}

	if token1 == token2 {
		t.Error("Tokens should be unique")
	}

	// 16 bytes = 32 hex chars
	if len(token1) != 32 {
		t.Errorf("Expected token length 32, got %d", len(token1))
	}
}
