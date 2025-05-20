package onntrackclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		if r.URL.Path != "/homepage/login" {
			t.Errorf("Expected request to '/homepage/login', got '%s'", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got '%s'", r.Method)
		}
		if r.Header.Get("Authorization") != "null" {
			t.Errorf("Expected Authorization header 'null', got '%s'", r.Header.Get("Authorization"))
		}
		if r.Header.Get("must") != "true" {
			t.Errorf("Expected must header 'true', got '%s'", r.Header.Get("must"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Decode request body
		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Check request body
		if loginReq.Account != "test@example.com" {
			t.Errorf("Expected account 'test@example.com', got '%s'", loginReq.Account)
		}
		if loginReq.Password != "password123" {
			t.Errorf("Expected password 'password123', got '%s'", loginReq.Password)
		}
		if loginReq.Language != "en" {
			t.Errorf("Expected language 'en', got '%s'", loginReq.Language)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"token":   "test-token-123",
			"message": "Login successful",
		})
	}))
	defer server.Close()

	// Create a client that uses the test server
	client, _ := NewClient(
		"", // No API key needed for login
		WithBaseURL(server.URL),
	)

	// Create login request
	loginReq := &LoginRequest{
		Account:   "test@example.com",
		Password:  "password123",
		Language:  "en",
		ValidCode: "",
		NodeID:    "",
	}

	// Call the Login method
	loginResp, resp, err := client.Auth.Login(context.Background(), loginReq)

	// Check for errors
	if err != nil {
		t.Fatalf("Login returned unexpected error: %v", err)
	}

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Login returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Check parsed data
	if !loginResp.Success {
		t.Errorf("Login success = %v, want %v", loginResp.Success, true)
	}
	if loginResp.Token != "test-token-123" {
		t.Errorf("Login token = %v, want %v", loginResp.Token, "test-token-123")
	}
	if loginResp.Message != "Login successful" {
		t.Errorf("Login message = %v, want %v", loginResp.Message, "Login successful")
	}

	// Check that the client's API key was updated
	if client.APIKey != "test-token-123" {
		t.Errorf("Client APIKey = %v, want %v", client.APIKey, "test-token-123")
	}
}
