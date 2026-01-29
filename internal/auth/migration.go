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
	tokenPath, errToken := config.GetTokenPath()
	if errToken != nil {
		return nil // Ignore error here, we'll hit it elsewhere if it's real
	}

	// 1. Check if old token.json exists
	if _, errStat := os.Stat(tokenPath); os.IsNotExist(errStat) {
		return nil // No old token, no migration needed
	}

	accountsDir, errAcc := config.GetAccountsDir()
	if errAcc != nil {
		return nil
	}

	// 2. Check if accounts directory already exists and is not empty
	if entries, errRead := os.ReadDir(accountsDir); errRead == nil && len(entries) > 0 {
		// Accounts already exist, assume migration happened or user already used multi-account
		return nil
	}

	fmt.Println(color.CyanString("🔄 Migrating from single account to multi-account format..."))

	// 3. Load old token
	data, errFile := os.ReadFile(tokenPath)
	if errFile != nil {
		return fmt.Errorf("failed to read old token for migration: %w", errFile)
	}

	var token TokenData
	if errUnmarshal := json.Unmarshal(data, &token); errUnmarshal != nil {
		return fmt.Errorf("failed to parse old token for migration: %w", errUnmarshal)
	}

	// 4. Ensure we have an email
	email := token.Email
	if email == "" {
		// If email is missing, we try to fetch it using the access token
		var errFetch error
		email, errFetch = fetchUserEmail(context.Background(), token.AccessToken)
		if errFetch != nil {
			return fmt.Errorf("failed to fetch user email during migration: %w. please run 'ag-quota login' instead", errFetch)
		}
		token.Email = email
	}

	// 5. Migrate to accounts/{email}.json
	if errSave := SaveTokenForAccount(email, &token); errSave != nil {
		return fmt.Errorf("failed to save migrated token: %w", errSave)
	}

	// 6. Set as default in config.json
	mgr, errMgr := NewAccountManager()
	if errMgr != nil {
		return fmt.Errorf("failed to initialize account manager during migration: %w", errMgr)
	}

	if errDef := mgr.SetDefaultAccount(email); errDef != nil {
		return fmt.Errorf("failed to set default account during migration: %w", errDef)
	}

	// 7. Backup old token.json
	backupPath := tokenPath + ".bak"
	if errRen := os.Rename(tokenPath, backupPath); errRen != nil {
		return fmt.Errorf("failed to backup old token: %w", errRen)
	}

	fmt.Println(color.GreenString("✅ Migration complete! Account %s is now your default.", email))
	subtle := color.New(color.Faint).Sprintf("   Note: Your old token has been backed up to %s", filepath.Base(backupPath))
	fmt.Println(subtle)
	fmt.Println()

	return nil
}
