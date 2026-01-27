package auth

import (
	"os"
	"testing"

	"github.com/gundamkid/anti-gravity-quota/internal/config"
)

func TestAccountManager(t *testing.T) {
	// Setup: Sử dụng thư mục tạm để không ảnh hưởng đến máy thật
	tmpDir := t.TempDir()
	originalConfigHome := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", originalConfigHome)

	mgr, err := NewAccountManager()
	if err != nil {
		t.Fatalf("Failed to create AccountManager: %v", err)
	}

	email := "test@example.com"
	token := &TokenData{
		AccessToken: "fake-access-token",
		Email:       email,
	}

	t.Run("SaveAccount", func(t *testing.T) {
		err := SaveTokenForAccount(email, token)
		if err != nil {
			t.Errorf("SaveTokenForAccount failed: %v", err)
		}

		// Kiểm tra file có tồn tại không
		path, _ := config.GetAccountPath(email)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Account file was not created at %s", path)
		}
	})

	t.Run("SetDefaultAccount", func(t *testing.T) {
		err := mgr.SetDefaultAccount(email)
		if err != nil {
			t.Errorf("SetDefaultAccount failed: %v", err)
		}

		cfg, _ := mgr.LoadConfig()
		if cfg.DefaultAccount != email {
			t.Errorf("Expected default account %s, got %s", email, cfg.DefaultAccount)
		}
	})

	t.Run("ListAccounts", func(t *testing.T) {
		accounts, err := mgr.ListAccounts()
		if err != nil {
			t.Errorf("ListAccounts failed: %v", err)
		}

		if len(accounts) != 1 {
			t.Errorf("Expected 1 account, got %d", len(accounts))
		}

		if accounts[0].Email != email {
			t.Errorf("Expected email %s, got %s", email, accounts[0].Email)
		}

		if !accounts[0].IsDefault {
			t.Error("Expected account to be marked as default")
		}
	})
}
