package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"golang.org/x/oauth2"
)

var (
	// locks is a map of mutexes per account to ensure thread safety
	locks sync.Map
)

func getLock(email string) *sync.Mutex {
	val, _ := locks.LoadOrStore(email, &sync.Mutex{})
	lock, _ := val.(*sync.Mutex)
	return lock
}

// TokenData represents stored OAuth2 token information
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
	Email        string    `json:"email,omitempty"`
	TierName     string    `json:"tier_name,omitempty"`
}

// SaveToken saves the token data to the current default account
func SaveToken(token *TokenData) error {
	mgr, err := NewAccountManager()
	if err != nil {
		return err
	}

	cfg, err := mgr.LoadConfig()
	if err != nil {
		return err
	}

	email := token.Email
	if email == "" {
		email = cfg.DefaultAccount
	}

	if email == "" {
		// If no email provided and no default, we can't save in the new format yet
		// This should only happen during initial login or before migration
		// For now, let's just use the old fallback if it exists, or error
		return fmt.Errorf("no email provided and no default account set")
	}

	return SaveTokenForAccount(email, token)
}

// LoadToken loads the token data from the current default account
func LoadToken() (*TokenData, error) {
	mgr, err := NewAccountManager()
	if err != nil {
		return nil, err
	}

	cfg, err := mgr.LoadConfig()
	if err != nil {
		return nil, err
	}

	if cfg.DefaultAccount == "" {
		return nil, fmt.Errorf("no default account set. please login first")
	}

	return LoadTokenForAccount(cfg.DefaultAccount)
}

// SaveTokenForAccount saves token data for a specific account
func SaveTokenForAccount(email string, token *TokenData) error {
	lock := getLock(email)
	lock.Lock()
	defer lock.Unlock()

	return saveTokenForAccount(email, token)
}

// saveTokenForAccount is the internal version of SaveTokenForAccount that doesn't acquire locks.
// It must be called while holding the lock for the account.
func saveTokenForAccount(email string, token *TokenData) error {
	if _, err := config.EnsureAccountsDir(); err != nil {
		return err
	}

	tokenPath, err := config.GetAccountPath(email)
	if err != nil {
		return err
	}

	// Ensure email is set in token
	token.Email = email

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	return config.AtomicWrite(tokenPath, data, 0600)
}

// LoadTokenForAccount loads token data for a specific account
func LoadTokenForAccount(email string) (*TokenData, error) {
	lock := getLock(email)
	lock.Lock()
	defer lock.Unlock()

	return loadTokenForAccount(email)
}

// loadTokenForAccount is the internal version of LoadTokenForAccount that doesn't acquire locks.
// It must be called while holding the lock for the account.
func loadTokenForAccount(email string) (*TokenData, error) {
	tokenPath, err := config.GetAccountPath(email)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("token for account %s not found", email)
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var token TokenData
	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &token, nil
}

// DeleteToken removes the stored token file for the current default account
func DeleteToken() error {
	mgr, err := NewAccountManager()
	if err != nil {
		return err
	}

	cfg, err := mgr.LoadConfig()
	if err != nil {
		return err
	}

	if cfg.DefaultAccount == "" {
		return nil // Nothing to delete
	}

	tokenPath, err := config.GetAccountPath(cfg.DefaultAccount)
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

// GetValidToken returns a valid access token for the default account, refreshing if necessary
func GetValidToken(oauthConfig *oauth2.Config) (string, error) {
	token, err := LoadToken()
	if err != nil {
		return "", err
	}

	return GetValidTokenForAccount(token.Email, oauthConfig)
}

// GetValidTokenForAccount returns a valid access token for a specific account, refreshing if necessary
func GetValidTokenForAccount(email string, oauthConfig *oauth2.Config) (string, error) {
	// Use a lock for this account to prevent concurrent refreshes
	lock := getLock(email)
	lock.Lock()
	defer lock.Unlock()

	// Load existing token for the account again while holding the lock
	// to ensure we have the latest one (it might have been refreshed by another goroutine)
	token, err := loadTokenForAccount(email)
	if err != nil {
		return "", err
	}

	// Check if token is still valid (another goroutine might have refreshed it)
	if token.IsValid() {
		return token.AccessToken, nil
	}

	// Token expired, try to refresh
	if token.RefreshToken == "" {
		return "", fmt.Errorf("token expired and no refresh token available for %s, please login again", email)
	}

	// Create token source
	tokenSource := oauthConfig.TokenSource(context.Background(), token.ToOAuth2Token())

	// Get fresh token
	newToken, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to refresh token for %s: %w", email, err)
	}

	// Create new TokenData with refreshed token
	refreshedToken := FromOAuth2Token(newToken, email)

	// Save the refreshed token for this account
	if err := saveTokenForAccount(email, refreshedToken); err != nil {
		return "", fmt.Errorf("failed to save refreshed token for %s: %w", email, err)
	}

	return refreshedToken.AccessToken, nil
}
