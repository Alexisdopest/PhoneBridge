package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/clipboard", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ClipboardHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("ClipboardHandler GET returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	req2, _ := http.NewRequest("GET", "/api/upload", nil)
	rr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(UploadHandler)
	handler2.ServeHTTP(rr2, req2)
	
	if status := rr2.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("UploadHandler GET returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
