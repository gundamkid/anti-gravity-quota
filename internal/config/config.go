package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application's configuration.
type Config struct {
	DefaultAccount string               `json:"default_account,omitempty"`
	Notifications  NotificationSettings `json:"notifications,omitempty"`
}

// NotificationSettings contains settings for various notification channels.
type NotificationSettings struct {
	Enabled  bool             `json:"enabled"`
	Telegram TelegramSettings `json:"telegram,omitempty"`
}

// TelegramSettings contains credentials for Telegram bot notifications.
type TelegramSettings struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

// AtomicWrite writes data to a file atomically by writing to a temp file first and then renaming it.
func AtomicWrite(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	// Clean up on error
	defer func() {
		if err != nil {
			os.Remove(tmpName)
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	if err = os.Chmod(tmpName, perm); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
}

const (
	AppName        = "ag-quota"
	TokenFileName  = "token.json" // Deprecated: use accounts/{email}.json
	ConfigFileName = "config.json"
	AccountsDir    = "accounts"
)

// GetAccountsDir returns the directory where account tokens are stored
func GetAccountsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, AccountsDir), nil
}

// EnsureAccountsDir creates the accounts directory if it doesn't exist
func EnsureAccountsDir() (string, error) {
	accountsDir, err := GetAccountsDir()
	if err != nil {
		return "", err
	}

	// Create directory with 0700 permissions (owner only)
	err = os.MkdirAll(accountsDir, 0700)
	if err != nil {
		return "", err
	}

	return accountsDir, nil
}

// GetAccountPath returns the full path to an account token file
func GetAccountPath(email string) (string, error) {
	accountsDir, err := GetAccountsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(accountsDir, email+".json"), nil
}

// GetConfigDir returns the configuration directory path
// Uses XDG_CONFIG_HOME on Linux/Mac, falls back to ~/.config/ag-quota
func GetConfigDir() (string, error) {
	var configDir string

	// Check for XDG_CONFIG_HOME environment variable
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		configDir = filepath.Join(xdg, AppName)
	} else {
		// Fall back to ~/.config/ag-quota
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config", AppName)
	}

	return configDir, nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	// Create directory with 0700 permissions (owner only)
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return "", err
	}

	return configDir, nil
}

// GetTokenPath returns the full path to the token file
func GetTokenPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, TokenFileName), nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ConfigFileName), nil
}

// LoadConfig loads the application configuration from the default path.
func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves the application configuration to the default path.
func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	if _, err := EnsureConfigDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return AtomicWrite(path, data, 0600)
}
