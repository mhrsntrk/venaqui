package realdebrid

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.apiToken != "test-token" {
		t.Errorf("NewClient() apiToken = %v, want test-token", client.apiToken)
	}
	if client.httpClient == nil {
		t.Error("NewClient() httpClient is nil")
	}
}

func TestUnrestrictLink_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/unrestrict/link" {
			t.Errorf("Expected /unrestrict/link, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Bearer test-token, got %s", r.Header.Get("Authorization"))
		}

		// Return success response
		response := UnrestrictedLink{
			ID:       "test-id",
			Filename: "test.zip",
			Filesize: 1024,
			Link:     "https://direct-download.com/file.zip",
			Host:     "example.com",
			Chunks:   1,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClientWithBaseURL("test-token", server.URL)

	// Test unrestrict
	result, err := client.UnrestrictLink("https://example.com/file.zip")
	if err != nil {
		t.Fatalf("UnrestrictLink() error = %v", err)
	}

	if result.ID != "test-id" {
		t.Errorf("UnrestrictLink() ID = %v, want test-id", result.ID)
	}
	if result.Filename != "test.zip" {
		t.Errorf("UnrestrictLink() Filename = %v, want test.zip", result.Filename)
	}
	if result.Link != "https://direct-download.com/file.zip" {
		t.Errorf("UnrestrictLink() Link = %v, want https://direct-download.com/file.zip", result.Link)
	}
}

func TestUnrestrictLink_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := ErrorResponse{
			Error:     "Invalid link",
			ErrorCode: 400,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
	}))
	defer server.Close()

	client := NewClientWithBaseURL("test-token", server.URL)

	_, err := client.UnrestrictLink("https://invalid.com/file.zip")
	if err == nil {
		t.Fatal("UnrestrictLink() expected error, got nil")
	}
}

func TestUnrestrictLink_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClientWithBaseURL("invalid-token", server.URL)

	_, err := client.UnrestrictLink("https://example.com/file.zip")
	if err == nil {
		t.Fatal("UnrestrictLink() expected error for unauthorized, got nil")
	}
}
