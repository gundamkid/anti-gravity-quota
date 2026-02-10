package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

// DisplayOptions controls how quota information is displayed
type DisplayOptions struct {
	Compact bool
}

// ClearTerminal clears the terminal screen
func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

// DisplayQuotaSummaryJSON displays quota information in JSON format
func DisplayQuotaSummaryJSON(summary *models.QuotaSummary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// DisplayQuotaSummary displays quota information in a formatted table
func DisplayQuotaSummary(summary *models.QuotaSummary, opts DisplayOptions) {
	// Header
	fmt.Println()
	color.Cyan("  ✨ Anti-Gravity Quota Status")
	if summary.Email != "" {
		tier := summary.TierName
		if tier == "" {
			tier = "Free 📦"
		}
		fmt.Printf("  📧 %s [%s]\n", summary.Email, tier)
	}
	fmt.Println()

	fmt.Printf("  Fetched: %s\n", summary.FetchedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Println()

	// Sort models by display name
	models := make([]models.ModelQuota, len(summary.Models))
	copy(models, summary.Models)
	sort.Slice(models, func(i, j int) bool {
		return models[i].DisplayName < models[j].DisplayName
	})

	// Create table
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

	if opts.Compact {
		t.AppendHeader(table.Row{"Model", "Quota", "Reset In"})
	} else {
		t.AppendHeader(table.Row{"Model", "Quota", "Reset In", "Status"})
	}

	for _, model := range models {
		if model.DisplayName == "" {
			continue
		}

		displayName := model.DisplayName
		if opts.Compact {
			displayName = shortenModelName(displayName)
		}

		percentage := model.GetRemainingPercentage()

		// Colorize Quota cell
		var quotaColor text.Colors
		if percentage <= 0 {
			quotaColor = text.Colors{text.FgHiBlack}
		} else if percentage <= 20 {
			quotaColor = text.Colors{text.FgRed, text.Bold}
		} else if percentage <= 50 {
			quotaColor = text.Colors{text.FgYellow}
		} else {
			quotaColor = text.Colors{text.FgGreen}
		}
		quotaStr := fmt.Sprintf("%3d%%", percentage)

		// Format Status with colors
		statusStr := model.GetStatusString()
		var statusColor text.Colors
		switch statusStr {
		case "HEALTHY":
			statusStr = "✓ HEALTHY"
			statusColor = text.Colors{text.FgGreen}
		case "WARNING":
			statusStr = "⚠ WARNING"
			statusColor = text.Colors{text.FgYellow}
		case "CRITICAL":
			statusStr = "⚡ CRITICAL"
			statusColor = text.Colors{text.FgRed}
		case "EMPTY":
			statusStr = "✗ EMPTY"
			statusColor = text.Colors{text.FgHiBlack}
		}

		if opts.Compact {
			t.AppendRow(table.Row{
				displayName,
				quotaColor.Sprint(quotaStr),
				formatResetTime(model, summary.FetchedAt),
			})
		} else {
			t.AppendRow(table.Row{
				displayName,
				quotaColor.Sprint(quotaStr),
				formatResetTime(model, summary.FetchedAt),
				statusColor.Sprint(statusStr),
			})
		}
	}

	// Indent the table slightly for better look
	// t.Render() returns the string. We prepend "  " to each line.
	rendered := t.Render()
	indented := "  " + strings.ReplaceAll(rendered, "\n", "\n  ")
	fmt.Println(indented)
	fmt.Println()

	// Footer with default model
	if summary.DefaultModelID != "" {
		for _, model := range models {
			if model.ModelID == summary.DefaultModelID {
				displayName := model.DisplayName
				if opts.Compact {
					displayName = shortenModelName(displayName)
				}
				color.Cyan("  ⭐ Default Model: %s", displayName)
				break
			}
		}
		fmt.Println()
	}
}

// formatResetTime formats the time until reset in a human-readable format
func formatResetTime(model models.ModelQuota, now time.Time) string {
	duration := model.ResetTime.Sub(now)

	if duration < 0 {
		return "Regenerating..."
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %02dh", days, hours)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh %02dm", hours, minutes)
	}

	return fmt.Sprintf("%dm", minutes)
}

// DisplayError displays an error message
func DisplayError(message string, err error) {
	color.Red("Error: %s", message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %v\n", err)
	}
}

// DisplayNotLoggedIn displays a message when user is not logged in
func DisplayNotLoggedIn() {
	color.Red("Not logged in")
	fmt.Println()
	fmt.Println("Please run the following command to authenticate:")
	color.Cyan("  ag-quota login")
	fmt.Println()
}

// DisplayLoading displays a loading message
func DisplayLoading(message string) {
	fmt.Printf("%s", message)
}

// DisplaySuccess displays a success message
func DisplaySuccess(message string) {
	color.Green("✓ %s", message)
}

// Spinner represents a simple text spinner
type Spinner struct {
	frames []string
	index  int
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:  0,
	}
}

// Next returns the next spinner frame
func (s *Spinner) Next() string {
	frame := s.frames[s.index]
	s.index = (s.index + 1) % len(s.frames)
	return frame
}

// AccountQuotaResult represents quota information for a single account
type AccountQuotaResult struct {
	Email        string               `json:"email"`
	QuotaSummary *models.QuotaSummary `json:"quota_summary,omitempty"`
	Error        string               `json:"error,omitempty"`
}

// DisplayAllAccountsQuotaJSON displays quota for all accounts in JSON format
func DisplayAllAccountsQuotaJSON(results []*AccountQuotaResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// DisplayAllAccountsQuota displays quota for all accounts in a formatted table
func DisplayAllAccountsQuota(results []*AccountQuotaResult, opts DisplayOptions) {
	if len(results) == 0 {
		color.Yellow("No accounts to display")
		return
	}

	// Header
	fmt.Println()
	color.Cyan("  ✨ Anti-Gravity Quota Status - All Accounts")
	fmt.Println()

	// Create a table for each account
	for _, result := range results {
		// Account header
		if result.Error != "" {
			color.Red("  ✗ %s", result.Email)
			fmt.Printf("    Error: %s\n", result.Error)
			fmt.Println()
			continue
		}

		if result.QuotaSummary == nil {
			color.Yellow("  ⚠ %s - No data", result.Email)
			fmt.Println()
			continue
		}

		// Display account email and tier
		tier := result.QuotaSummary.TierName
		if tier == "" {
			tier = "Free 📦"
		}
		color.Cyan("  📧 %s [%s]", result.Email, tier)
		fmt.Println()

		// Sort models by display name
		models := make([]models.ModelQuota, len(result.QuotaSummary.Models))
		copy(models, result.QuotaSummary.Models)
		sort.Slice(models, func(i, j int) bool {
			return models[i].DisplayName < models[j].DisplayName
		})

		// Create table
		t := table.NewWriter()
		t.SetStyle(table.StyleRounded)

		// Customize style
		style := table.StyleRounded
		style.Color.Header = text.Colors{text.FgCyan, text.Bold}
		style.Color.Border = text.Colors{text.FgCyan}
		style.Color.Separator = text.Colors{text.FgCyan}
		t.SetStyle(style)

		if opts.Compact {
			t.AppendHeader(table.Row{"Model", "Quota", "Reset In"})
		} else {
			t.AppendHeader(table.Row{"Model", "Quota", "Reset In", "Status"})
		}

		for _, model := range models {
			if model.DisplayName == "" {
				continue
			}

			displayName := model.DisplayName
			if opts.Compact {
				displayName = shortenModelName(displayName)
			}

			percentage := model.GetRemainingPercentage()

			// Colorize Quota cell
			var quotaColor text.Colors
			if percentage <= 0 {
				quotaColor = text.Colors{text.FgHiBlack}
			} else if percentage <= 20 {
				quotaColor = text.Colors{text.FgRed, text.Bold}
			} else if percentage <= 50 {
				quotaColor = text.Colors{text.FgYellow}
			} else {
				quotaColor = text.Colors{text.FgGreen}
			}
			quotaStr := fmt.Sprintf("%3d%%", percentage)

			// Format Status with colors
			statusStr := model.GetStatusString()
			var statusColor text.Colors
			switch statusStr {
			case "HEALTHY":
				statusStr = "✓ HEALTHY"
				statusColor = text.Colors{text.FgGreen}
			case "WARNING":
				statusStr = "⚠ WARNING"
				statusColor = text.Colors{text.FgYellow}
			case "CRITICAL":
				statusStr = "⚡ CRITICAL"
				statusColor = text.Colors{text.FgRed}
			case "EMPTY":
				statusStr = "✗ EMPTY"
				statusColor = text.Colors{text.FgHiBlack}
			}

			if opts.Compact {
				t.AppendRow(table.Row{
					displayName,
					quotaColor.Sprint(quotaStr),
					formatResetTime(model, result.QuotaSummary.FetchedAt),
				})
			} else {
				t.AppendRow(table.Row{
					displayName,
					quotaColor.Sprint(quotaStr),
					formatResetTime(model, result.QuotaSummary.FetchedAt),
					statusColor.Sprint(statusStr),
				})
			}
		}

		// Indent the table
		rendered := t.Render()
		indented := "    " + strings.ReplaceAll(rendered, "\n", "\n    ")
		fmt.Println(indented)
		fmt.Println()
	}
}

// shortenModelName reduces model name length for compact mode
func shortenModelName(name string) string {
	name = strings.ReplaceAll(name, "Claude ", "")
	name = strings.ReplaceAll(name, "Gemini ", "Gem ")
	name = strings.ReplaceAll(name, "(Thinking)", "(T)")
	name = strings.ReplaceAll(name, "(thinking)", "(T)")
	name = strings.ReplaceAll(name, "(Low)", "(L)")
	name = strings.ReplaceAll(name, "(Medium)", "(M)")
	name = strings.ReplaceAll(name, "(High)", "(H)")
	return name
}

// GetTerminalWidth returns the width of the terminal
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // Default fallback
	}
	return width
}

// DisplayWatchHeader displays the header for watch mode
func DisplayWatchHeader(interval int) {
	ClearTerminal()
	fmt.Println()
	color.HiCyan("  👀 Anti-Gravity Quota - WATCH MODE ACTIVE")
	color.HiBlack("  Interval: %dm | Auto-refreshing", interval)
	fmt.Println()
}

// DisplayWatchFooter displays the footer for watch mode
func DisplayWatchFooter(lastUpdated time.Time) {
	fmt.Println()
	color.HiBlack("  Refreshed at: %s", lastUpdated.Format("15:04:05"))
	color.HiBlack("  Press Ctrl+C to exit")
}
