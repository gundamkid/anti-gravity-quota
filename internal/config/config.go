package config

import (
	"os"
	"path/filepath"
)

const (
	AppName        = "ag-quota"
	TokenFileName  = "token.json"
	ConfigFileName = "config.json"
)

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
