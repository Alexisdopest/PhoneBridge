package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	validToken := "123456"
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	handler := Middleware(validToken, nextHandler)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{"Valid Token", "Bearer 123456", http.StatusOK},
		{"Invalid Token", "Bearer wrong", http.StatusUnauthorized},
		{"Missing Bearer", "123456", http.StatusUnauthorized},
		{"Missing Header", "", http.StatusUnauthorized},
		{"Empty Token", "Bearer ", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
