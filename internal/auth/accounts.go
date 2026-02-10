package auth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/config"
)

var (
	// ErrAccountNotFound is returned when an account cannot be found
	ErrAccountNotFound = errors.New("account not found")
)

// AccountInfo represents metadata about a saved account
type AccountInfo struct {
	Email      string    `json:"email"`
	TierName   string    `json:"tier_name"`
	IsDefault  bool      `json:"is_default"`
	LastUsed   time.Time `json:"last_used"`
	TokenValid bool      `json:"token_valid"`
}

// AccountManager handles all account-related operations
type AccountManager struct {
	accountsDir string
	configPath  string
}

// NewAccountManager creates a new instance of AccountManager
func NewAccountManager() (*AccountManager, error) {
	accountsDir, err := config.EnsureAccountsDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup accounts directory: %w", err)
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	return &AccountManager{
		accountsDir: accountsDir,
		configPath:  configPath,
	}, nil
}

// ListAccounts returns a list of all saved accounts
func (m *AccountManager) ListAccounts() ([]AccountInfo, error) {
	entries, err := os.ReadDir(m.accountsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts directory: %w", err)
	}

	appCfg, err := m.LoadConfig()
	if err != nil {
		return nil, err
	}

	var accounts []AccountInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		email := entry.Name()[:len(entry.Name())-len(".json")]

		// Load token to check validity and get tier
		token, err := LoadTokenForAccount(email)
		valid := err == nil && token.IsValid()
		tier := ""
		if err == nil {
			tier = token.TierName
		}

		accounts = append(accounts, AccountInfo{
			Email:      email,
			TierName:   tier,
			IsDefault:  email == appCfg.DefaultAccount,
			TokenValid: valid,
		})
	}

	// Sort: default first, then alphabetical
	sort.Slice(accounts, func(i, j int) bool {
		if accounts[i].IsDefault {
			return true
		}
		if accounts[j].IsDefault {
			return false
		}
		return accounts[i].Email < accounts[j].Email
	})

	return accounts, nil
}

// LoadConfig loads the application config
func (m *AccountManager) LoadConfig() (*config.Config, error) {
	return config.LoadConfig()
}

// SaveConfig saves the application config
func (m *AccountManager) SaveConfig(cfg *config.Config) error {
	return config.SaveConfig(cfg)
}

// SetDefaultAccount sets the default account email
func (m *AccountManager) SetDefaultAccount(email string) error {
	// Verify account exists
	path, err := config.GetAccountPath(email)
	if err != nil {
		return err
	}
	if _, errStat := os.Stat(path); errStat != nil {
		if os.IsNotExist(errStat) {
			return fmt.Errorf("%w: %s", ErrAccountNotFound, email)
		}
		return errStat
	}

	cfg, err := m.LoadConfig()
	if err != nil {
		return err
	}

	cfg.DefaultAccount = email
	return m.SaveConfig(cfg)
}

// RemoveAccount deletes an account's token and clears it from config if it was default
func (m *AccountManager) RemoveAccount(email string) error {
	path, err := config.GetAccountPath(email)
	if err != nil {
		return err
	}

	// Remove token file
	if errRemove := os.Remove(path); errRemove != nil {
		if os.IsNotExist(errRemove) {
			return fmt.Errorf("%w: %s", ErrAccountNotFound, email)
		}
		return fmt.Errorf("failed to delete account token: %w", errRemove)
	}

	// Update config if it was the default account
	cfg, err := m.LoadConfig()
	if err != nil {
		return err
	}

	if cfg.DefaultAccount == email {
		cfg.DefaultAccount = ""
		return m.SaveConfig(cfg)
	}

	return nil
}
