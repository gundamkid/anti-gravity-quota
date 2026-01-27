package auth

import (
	"encoding/base64"
	"testing"
)

func TestGenerateCodeVerifier(t *testing.T) {
	v1, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("GenerateCodeVerifier failed: %v", err)
	}

	if len(v1) == 0 {
		t.Error("verifier length is 0")
	}

	v2, _ := GenerateCodeVerifier()
	if v1 == v2 {
		t.Error("GenerateCodeVerifier returned identical strings")
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "test-verifier-string"
	challenge := GenerateCodeChallenge(verifier)

	if len(challenge) == 0 {
		t.Error("challenge length is 0")
	}

	// Verify it's valid base64url
	_, err := base64.RawURLEncoding.DecodeString(challenge)
	if err != nil {
		t.Errorf("GenerateCodeChallenge returned invalid base64url: %v", err)
	}
}

func TestGenerateState(t *testing.T) {
	s1, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState failed: %v", err)
	}

	if len(s1) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("expected length 32, got %d", len(s1))
	}
}
