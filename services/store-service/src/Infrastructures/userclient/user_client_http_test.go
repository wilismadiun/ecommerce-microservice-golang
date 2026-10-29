package userclient

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserCLientHttp_IsUserExists(t *testing.T) {
	// Mock Server untuk mensimulasikan User Service
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/user-123/profile" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/users/user-404/profile" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Skenario status code tidak terduga
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Set environment variable USERSERVICEURL agar mengarah ke Mock Server
	os.Setenv("USERSERVICEURL", server.URL)
	defer os.Unsetenv("USERSERVICEURL")

	client := &UserCLientHttp{}

	t.Run("User Not Found (404)", func(t *testing.T) {
		exists, err := client.IsUserExists("user-404")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if exists {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("Unexpected Status Code (500)", func(t *testing.T) {
		_, err := client.IsUserExists("user-500")
		if err == nil {
			t.Errorf("expected an error, got nil")
		}
	})

	t.Run("User Exists (200)", func(t *testing.T) {
		exists, err := client.IsUserExists("user-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !exists {
			t.Errorf("expected true, got false")
		}
	})
}
