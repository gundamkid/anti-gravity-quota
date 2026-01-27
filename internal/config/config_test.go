package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	// Test with XDG_CONFIG_HOME
	expectedXDG := "/tmp/ag-quota-test-xdg"
	os.Setenv("XDG_CONFIG_HOME", expectedXDG)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}
	expected := filepath.Join(expectedXDG, AppName)
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}

	// Test without XDG_CONFIG_HOME
	os.Unsetenv("XDG_CONFIG_HOME")
	home, _ := os.UserHomeDir()
	dir, err = GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}
	expected = filepath.Join(home, ".config", AppName)
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	dir, err := EnsureConfigDir()
	if err != nil {
		t.Fatalf("EnsureConfigDir failed: %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("directory %s was not created", dir)
	}
}

func TestGetAccountPath(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	email := "test@example.com"
	path, err := GetAccountPath(email)
	if err != nil {
		t.Fatalf("GetAccountPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, AppName, AccountsDir, email+".json")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-file.txt")
	content := []byte("hello world")

	err := AtomicWrite(path, content, 0644)
	if err != nil {
		t.Fatalf("AtomicWrite failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("expected %s, got %s", string(content), string(data))
	}
}
