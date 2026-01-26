package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateCodeVerifier generates a PKCE code verifier
// Returns a random string of 43-128 characters
func GenerateCodeVerifier() (string, error) {
	// Generate 32 random bytes (will be 64 hex characters)
	return generateRandomString(32)
}

// GenerateCodeChallenge generates a PKCE code challenge from a verifier
// Uses SHA256 hash and base64url encoding
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	// Use base64url encoding without padding
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// GenerateState generates a random state parameter for OAuth2
func GenerateState() (string, error) {
	return generateRandomString(16)
}
