package auth

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTokenData_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		token    TokenData
		expected bool
	}{
		{
			name: "Valid token",
			token: TokenData{
				AccessToken: "valid",
				Expiry:      time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Expired token",
			token: TokenData{
				AccessToken: "expired",
				Expiry:      time.Now().Add(-time.Hour),
			},
			expected: false,
		},
		{
			name: "Empty token",
			token: TokenData{
				AccessToken: "",
				Expiry:      time.Now().Add(time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.token.IsValid() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.token.IsValid())
			}
		})
	}
}

func TestSaveLoadTokenForAccount(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	email := "user@example.com"
	token := &TokenData{
		AccessToken: "test-token",
		Email:       email,
		Expiry:      time.Now().Add(time.Hour),
	}

	// Test Save
	err := SaveTokenForAccount(email, token)
	if err != nil {
		t.Fatalf("SaveTokenForAccount failed: %v", err)
	}

	// Test Load
	loaded, err := LoadTokenForAccount(email)
	if err != nil {
		t.Fatalf("LoadTokenForAccount failed: %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Errorf("expected %s, got %s", token.AccessToken, loaded.AccessToken)
	}
	if loaded.Email != email {
		t.Errorf("expected %s, got %s", email, loaded.Email)
	}
}

func TestSaveLoadTokenConcurrent(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	email := "concurrent@example.com"
	const iterations = 100

	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func(id int) {
			token := &TokenData{
				AccessToken: fmt.Sprintf("token-%d", id),
				Email:       email,
				Expiry:      time.Now().Add(time.Hour),
			}
			_ = SaveTokenForAccount(email, token)
			_, _ = LoadTokenForAccount(email)
			done <- true
		}(i)
	}

	for i := 0; i < iterations; i++ {
		<-done
	}
}
