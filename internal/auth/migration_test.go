package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateIfNeeded(t *testing.T) {
	// Setup: Create a temporary config directory
	tempDir, err := os.MkdirTemp("", "ag-quota-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override GetConfigDir for testing
	// NOTE: This usually requires a way to inject the config path,
	// but here we might need to use environment variables or a global variable if implemented.
	// Looking at internal/config/config.go, it uses XDG_CONFIG_HOME.
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	configDir := filepath.Join(tempDir, config.AppName)
	err = os.MkdirAll(configDir, 0700)
	require.NoError(t, err)

	// 1. Create an old token.json
	oldToken := TokenData{
		AccessToken:  "old-access-token",
		RefreshToken: "old-refresh-token",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(1 * time.Hour),
		Email:        "test-migrated@gmail.com",
	}
	tokenData, err := json.Marshal(oldToken)
	require.NoError(t, err)

	oldTokenPath := filepath.Join(configDir, "token.json")
	err = os.WriteFile(oldTokenPath, tokenData, 0600)
	require.NoError(t, err)

	// 2. Run migration
	err = MigrateIfNeeded()
	require.NoError(t, err)

	// 3. Verify migration results
	// Check new token file
	newAccountPath := filepath.Join(configDir, "accounts", "test-migrated@gmail.com.json")
	_, err = os.Stat(newAccountPath)
	assert.NoError(t, err, "Migrated account token should exist")

	// Check config.json
	configPath := filepath.Join(configDir, "config.json")
	var cfg AppConfig
	cfgData, err := os.ReadFile(configPath)
	require.NoError(t, err)
	err = json.Unmarshal(cfgData, &cfg)
	require.NoError(t, err)
	assert.Equal(t, "test-migrated@gmail.com", cfg.DefaultAccount)

	// Check backup
	_, err = os.Stat(oldTokenPath)
	assert.Error(t, err, "Old token should be moved/renamed")
	assert.True(t, os.IsNotExist(err))

	backupPath := oldTokenPath + ".bak"
	_, err = os.Stat(backupPath)
	assert.NoError(t, err, "Backup file should exist")
}
