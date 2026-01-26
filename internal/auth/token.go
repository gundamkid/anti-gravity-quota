package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"golang.org/x/oauth2"
)

// TokenData represents stored OAuth2 token information
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
	Email        string    `json:"email,omitempty"`
}

// SaveToken saves the token data to a JSON file
func SaveToken(token *TokenData) error {
	// Ensure config directory exists
	if _, err := config.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Get token file path
	tokenPath, err := config.GetTokenPath()
	if err != nil {
		return fmt.Errorf("failed to get token path: %w", err)
	}

	// Marshal token to JSON
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write to file with 0600 permissions (owner only)
	err = os.WriteFile(tokenPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads the token data from the JSON file
func LoadToken() (*TokenData, error) {
	tokenPath, err := config.GetTokenPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get token path: %w", err)
	}

	// Read token file
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in: token file not found")
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	// Unmarshal token
	var token TokenData
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &token, nil
}

// DeleteToken removes the stored token file
func DeleteToken() error {
	tokenPath, err := config.GetTokenPath()
	if err != nil {
		return fmt.Errorf("failed to get token path: %w", err)
	}

	err = os.Remove(tokenPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete token file: %w", err)
	}

	return nil
}

// IsExpired checks if the token is expired
func (t *TokenData) IsExpired() bool {
	return time.Now().After(t.Expiry)
}

// IsValid checks if the token exists and is not expired
func (t *TokenData) IsValid() bool {
	return t.AccessToken != "" && !t.IsExpired()
}

// ToOAuth2Token converts TokenData to oauth2.Token
func (t *TokenData) ToOAuth2Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenType:    t.TokenType,
		Expiry:       t.Expiry,
	}
}

// FromOAuth2Token creates TokenData from oauth2.Token
func FromOAuth2Token(token *oauth2.Token, email string) *TokenData {
	return &TokenData{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
		Email:        email,
	}
}

// RefreshToken refreshes an expired token using the refresh token
func RefreshToken(token *TokenData, oauthConfig *oauth2.Config) (*TokenData, error) {
	if token.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	// Create token source
	tokenSource := oauthConfig.TokenSource(context.Background(), token.ToOAuth2Token())

	// Get fresh token
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Create new TokenData with refreshed token
	refreshedToken := FromOAuth2Token(newToken, token.Email)

	// Save the refreshed token
	if err := SaveToken(refreshedToken); err != nil {
		return nil, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return refreshedToken, nil
}

// GetValidToken returns a valid access token, refreshing if necessary
func GetValidToken(oauthConfig *oauth2.Config) (string, error) {
	// Load existing token
	token, err := LoadToken()
	if err != nil {
		return "", err
	}

	// Check if token is still valid
	if token.IsValid() {
		return token.AccessToken, nil
	}

	// Token expired, try to refresh
	if token.RefreshToken == "" {
		return "", fmt.Errorf("token expired and no refresh token available, please login again")
	}

	// Refresh the token
	refreshedToken, err := RefreshToken(token, oauthConfig)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return refreshedToken.AccessToken, nil
}
