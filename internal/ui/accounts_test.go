package ui

import (
	"testing"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/auth"
)

func TestDisplayAccountsList(t *testing.T) {
	tests := []struct {
		name     string
		accounts []auth.AccountInfo
	}{
		{
			name: "Single default account",
			accounts: []auth.AccountInfo{
				{
					Email:      "user@gmail.com",
					TierName:   "Pro 💎",
					IsDefault:  true,
					LastUsed:   time.Now(),
					TokenValid: true,
				},
			},
		},
		{
			name: "Multiple accounts with default",
			accounts: []auth.AccountInfo{
				{
					Email:      "default@gmail.com",
					TierName:   "Ultra 🚀",
					IsDefault:  true,
					LastUsed:   time.Now(),
					TokenValid: true,
				},
				{
					Email:      "second@gmail.com",
					TierName:   "Free 📦",
					IsDefault:  false,
					LastUsed:   time.Now().Add(-24 * time.Hour),
					TokenValid: true,
				},
			},
		},
		{
			name: "Account with expired token",
			accounts: []auth.AccountInfo{
				{
					Email:      "expired@gmail.com",
					IsDefault:  true,
					LastUsed:   time.Now(),
					TokenValid: false,
				},
			},
		},
		{
			name:     "Empty accounts list",
			accounts: []auth.AccountInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test just ensures the function doesn't panic
			// Visual output should be manually verified
			DisplayAccountsList(tt.accounts)
		})
	}
}
