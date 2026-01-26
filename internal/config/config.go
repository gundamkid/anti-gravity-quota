package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	AppName        = "anti-gravity-quota"
	ConfigFileName = "config.yaml"
)

type Config struct {
	Accounts []Account `yaml:"accounts"`
	Display  Display   `yaml:"display"`
}

type Account struct {
	Email  string `yaml:"email"`
	Active bool   `yaml:"active"`
}

type Display struct {
	Format  string `yaml:"format"`  // table, json, compact
	Refresh int    `yaml:"refresh"` // seconds
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Accounts: []Account{}, // Empty initially
		Display: Display{
			Format:  "table",
			Refresh: 60,
		},
	}
}

// LoadConfig reads the configuration from the config file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// SaveConfig writes the configuration to the config file
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil { // Secure permissions
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the absolute path to the config file
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}
	return filepath.Join(configDir, AppName, ConfigFileName), nil
}

// GetActiveAccount returns the currently active account, if any
func (c *Config) GetActiveAccount() *Account {
	for i := range c.Accounts {
		if c.Accounts[i].Active {
			return &c.Accounts[i]
		}
	}
	return nil
}

// SetActiveAccount sets the active account by email
func (c *Config) SetActiveAccount(email string) error {
	found := false
	for i := range c.Accounts {
		if c.Accounts[i].Email == email {
			c.Accounts[i].Active = true
			found = true
		} else {
			c.Accounts[i].Active = false
		}
	}
	
	if !found {
		return fmt.Errorf("account not found: %s", email)
	}
	return nil
}

// AddAccount adds a new account to the list
func (c *Config) AddAccount(email string) {
	// Check if already exists
	for _, acc := range c.Accounts {
		if acc.Email == email {
			return
		}
	}
	
	c.Accounts = append(c.Accounts, Account{
		Email:  email,
		Active: len(c.Accounts) == 0, // Make active if it's the first one
	})
}
