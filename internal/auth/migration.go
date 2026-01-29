package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/config"
)

// MigrateIfNeeded checks if a migration from single-account format (token.json)
// to multi-account format (accounts/email.json) is needed, and performs it if so.
func MigrateIfNeeded() error {
	tokenPath, err := config.GetTokenPath()
	if err != nil {
		return nil // Ignore error here, we'll hit it elsewhere if it's real
	}

	// 1. Check if old token.json exists
	if _, errStat := os.Stat(tokenPath); os.IsNotExist(errStat) {
		return nil // No old token, no migration needed
	}

	accountsDir, err := config.GetAccountsDir()
	if err != nil {
		return nil
	}

	// 2. Check if accounts directory already exists and is not empty
	if entries, errRead := os.ReadDir(accountsDir); errRead == nil && len(entries) > 0 {
		// Accounts already exist, assume migration happened or user already used multi-account
		return nil
	}

	fmt.Println(color.CyanString("🔄 Migrating from single account to multi-account format..."))

	// 3. Load old token
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to read old token for migration: %w", err)
	}

	var token TokenData
	if errUnmarshal := json.Unmarshal(data, &token); errUnmarshal != nil {
		return fmt.Errorf("failed to parse old token for migration: %w", errUnmarshal)
	}

	// 4. Ensure we have an email
	email := token.Email
	if email == "" {
		// If email is missing, we try to fetch it using the access token
		email, err = fetchUserEmail(context.Background(), token.AccessToken)
		if err != nil {
			return fmt.Errorf("failed to fetch user email during migration: %w. please run 'ag-quota login' instead", err)
		}
		token.Email = email
	}

	// 5. Migrate to accounts/{email}.json
	if err := SaveTokenForAccount(email, &token); err != nil {
		return fmt.Errorf("failed to save migrated token: %w", err)
	}

	// 6. Set as default in config.json
	mgr, err := NewAccountManager()
	if err != nil {
		return fmt.Errorf("failed to initialize account manager during migration: %w", err)
	}

	if err := mgr.SetDefaultAccount(email); err != nil {
		return fmt.Errorf("failed to set default account during migration: %w", err)
	}

	// 7. Backup old token.json
	backupPath := tokenPath + ".bak"
	if err := os.Rename(tokenPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup old token: %w", err)
	}

	fmt.Println(color.GreenString("✅ Migration complete! Account %s is now your default.", email))
	subtle := color.New(color.Faint).Sprintf("   Note: Your old token has been backed up to %s", filepath.Base(backupPath))
	fmt.Println(subtle)
	fmt.Println()

	return nil
}
