package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/jedib0t/go-pretty/v6/table"
)

// DisplayAccountsList displays a formatted table of all saved accounts
func DisplayAccountsList(accounts []auth.AccountInfo) {
	fmt.Println()
	fmt.Println("  📋 Saved Accounts")
	fmt.Println()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)

	// Header
	t.AppendHeader(table.Row{"", "Account", "Status"})

	// Rows
	for _, acc := range accounts {
		marker := "  "
		if acc.IsDefault {
			marker = color.YellowString("★")
		}

		status := color.GreenString("✓ Valid")
		if !acc.TokenValid {
			status = color.RedString("✗ Expired")
		}

		t.AppendRow(table.Row{marker, acc.Email, status})
	}

	t.Render()
	fmt.Println()

	// Show legend
	fmt.Println("  " + color.YellowString("★") + " = Default account")
	fmt.Println()
}
