package ui

import (
	"fmt"

	"strings"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// DisplayAccountsList displays a formatted table of all saved accounts
func DisplayAccountsList(accounts []auth.AccountInfo) {
	fmt.Println()
	fmt.Println("  📋 Saved Accounts")
	fmt.Println()

	t := table.NewWriter()
	// No OutputMirror, we will Render() to string to indent it manually

	// Set style (Rounded is modern and clean)
	t.SetStyle(table.StyleRounded)

	// Customize style for specific look
	style := table.StyleRounded
	style.Color.Header = text.Colors{text.FgCyan, text.Bold}
	style.Color.Border = text.Colors{text.FgCyan}
	style.Color.Separator = text.Colors{text.FgCyan}
	t.SetStyle(style)

	// Header
	t.AppendHeader(table.Row{"", "Account", "Plan", "Status"})

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

		t.AppendRow(table.Row{marker, acc.Email, acc.TierName, status})
	}

	// Indent the table slightly for better look
	rendered := t.Render()
	indented := "  " + strings.ReplaceAll(rendered, "\n", "\n  ")
	fmt.Println(indented)
	fmt.Println()

	// Show legend
	fmt.Println("  " + color.YellowString("★") + " = Default account")
	fmt.Println()
}
