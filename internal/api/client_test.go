package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

func TestClient_doRequest(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("User-Agent") != UserAgent {
			t.Errorf("expected User-Agent %s, got %s", UserAgent, r.Header.Get("User-Agent"))
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	data, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest failed: %v", err)
	}

	if string(data) != `{"status": "ok"}` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestExtractProjectId(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"String ID", "project-123", "project-123"},
		{"Object ID", map[string]interface{}{"id": "project-456"}, "project-456"},
		{"Nil input", nil, ""},
		{"Empty string", "", ""},
		{"Invalid map", map[string]interface{}{"foo": "bar"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractProjectId(tt.input)
			if got != tt.expected {
				t.Errorf("extractProjectId() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClient_LoadCodeAssist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.LoadCodeAssistResponse{
			ProjectID: "test-project",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "fake-token", // Set token to avoid EnsureAuthenticated logic which needs actual storage
	}

	resp, err := client.LoadCodeAssist()
	if err != nil {
		t.Fatalf("LoadCodeAssist failed: %v", err)
	}

	if resp.ProjectID != "test-project" {
		t.Errorf("expected test-project, got %s", resp.ProjectID)
	}
	if client.projectID != "test-project" {
		t.Errorf("client projectID not updated")
	}
}
