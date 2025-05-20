package onntrackclient

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if client.APIKey != "test-api-key" {
		t.Errorf("NewClient APIKey = %v, want %v", client.APIKey, "test-api-key")
	}

	if client.BaseURL.String() != DefaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", client.BaseURL.String(), DefaultBaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("NewClient HTTPClient is nil")
	}
}

func TestWithBaseURL(t *testing.T) {
	client, err := NewClient(
		"test-api-key",
		WithBaseURL("https://custom-api.example.com"),
	)
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if client.BaseURL.String() != "https://custom-api.example.com" {
		t.Errorf("WithBaseURL = %v, want %v", client.BaseURL.String(), "https://custom-api.example.com")
	}
}

func TestWithLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client, err := NewClient(
		"test-api-key",
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	// Check that the transport is a LoggingTransport
	transport, ok := client.HTTPClient.Transport.(*LoggingTransport)
	if !ok {
		t.Fatalf("WithLogger did not set LoggingTransport, got %T", client.HTTPClient.Transport)
	}

	if transport.Logger != logger {
		t.Errorf("WithLogger Logger = %v, want %v", transport.Logger, logger)
	}
}

func TestClient_NewRequest(t *testing.T) {
	client, _ := NewClient("test-api-key")

	inURL, outURL := "foo", DefaultBaseURL+"foo"
	inBody, outBody := &DeviceCreateRequest{Name: "test-device"}, `{"name":"test-device","type":"","imei":""}`+"\n"
	req, _ := client.NewRequest(context.Background(), http.MethodPost, inURL, inBody)

	// Test that the URL was correctly formed
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL = %v, want %v", inURL, got, want)
	}

	// Test that the body was correctly encoded
	body := make([]byte, req.ContentLength)
	req.Body.Read(body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body = %v, want %v", inBody, got, want)
	}

	// Test that the correct headers were set
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest() Content-Type = %v, want %v", got, want)
	}
	if got, want := req.Header.Get("Accept"), "application/json"; got != want {
		t.Errorf("NewRequest() Accept = %v, want %v", got, want)
	}
	if got, want := req.Header.Get("Authorization"), "Bearer test-api-key"; got != want {
		t.Errorf("NewRequest() Authorization = %v, want %v", got, want)
	}
}

func TestClient_Do(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/devices" {
			t.Errorf("Expected request to '/devices', got '%s'", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header 'Bearer test-api-key', got '%s'", r.Header.Get("Authorization"))
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]*Device{
			{
				ID:          "device-1",
				Name:        "Test Device 1",
				Type:        "tracker",
				IMEI:        "123456789012345",
				Status:      "active",
				LastUpdated: "2023-01-01T12:00:00Z",
			},
		})
	}))
	defer server.Close()

	// Create a client that uses the test server
	client, _ := NewClient(
		"test-api-key",
		WithBaseURL(server.URL),
	)

	// Make a request
	req, _ := client.NewRequest(context.Background(), http.MethodGet, "devices", nil)
	var devices []*Device
	resp, err := client.Do(req, &devices)

	// Check for errors
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Do returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Check parsed data
	if len(devices) != 1 {
		t.Fatalf("Do returned %d devices, want 1", len(devices))
	}
	if devices[0].ID != "device-1" {
		t.Errorf("Device ID = %v, want %v", devices[0].ID, "device-1")
	}
	if devices[0].Name != "Test Device 1" {
		t.Errorf("Device Name = %v, want %v", devices[0].Name, "Test Device 1")
	}
}

func TestDevicesService_List(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/devices" {
			t.Errorf("Expected request to '/devices', got '%s'", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", r.Method)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]*Device{
			{
				ID:          "device-1",
				Name:        "Test Device 1",
				Type:        "tracker",
				IMEI:        "123456789012345",
				Status:      "active",
				LastUpdated: "2023-01-01T12:00:00Z",
			},
			{
				ID:          "device-2",
				Name:        "Test Device 2",
				Type:        "tracker",
				IMEI:        "987654321098765",
				Status:      "inactive",
				LastUpdated: "2023-01-02T12:00:00Z",
			},
		})
	}))
	defer server.Close()

	// Create a client that uses the test server
	client, _ := NewClient(
		"test-api-key",
		WithBaseURL(server.URL),
	)

	// Call the List method
	devices, resp, err := client.Devices.List(context.Background(), nil)

	// Check for errors
	if err != nil {
		t.Fatalf("List returned unexpected error: %v", err)
	}

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("List returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Check parsed data
	if len(devices) != 2 {
		t.Fatalf("List returned %d devices, want 2", len(devices))
	}
	if devices[0].ID != "device-1" {
		t.Errorf("Device ID = %v, want %v", devices[0].ID, "device-1")
	}
	if devices[1].ID != "device-2" {
		t.Errorf("Device ID = %v, want %v", devices[1].ID, "device-2")
	}
}
