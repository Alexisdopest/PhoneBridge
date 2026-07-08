package auth

import (
	"net/http"
	"strings"
)

// Middleware returns an http.HandlerFunc that checks the Bearer token
func Middleware(validToken string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized: Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		if parts[1] != validToken {
			http.Error(w, "Unauthorized: Token mismatch", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed
		next(w, r)
	}
}
