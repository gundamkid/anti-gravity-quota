package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

// DisplayQuotaSummary displays quota information in a formatted table
func DisplayQuotaSummary(summary *models.QuotaSummary) {
	// Header
	fmt.Println()
	color.Cyan("═══════════════════════════════════════════════════════════════")
	color.Cyan("              Anti-Gravity Quota Status")
	color.Cyan("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Account info
	if summary.Email != "" {
		fmt.Printf("  Account: %s\n", color.GreenString(summary.Email))
	}
	if summary.ProjectID != "" {
		fmt.Printf("  Project: %s\n", summary.ProjectID)
	}
	fmt.Printf("  Fetched: %s\n", summary.FetchedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Println()

	// Sort models by display name
	models := make([]models.ModelQuota, len(summary.Models))
	copy(models, summary.Models)
	sort.Slice(models, func(i, j int) bool {
		return models[i].DisplayName < models[j].DisplayName
	})

	// Create simple table using custom formatting
	printTableHeader()

	for _, model := range models {
		quotaStr := formatQuota(model)
		resetStr := formatResetTime(model)
		statusStr := formatStatus(model)

		printTableRow(model.DisplayName, quotaStr, resetStr, statusStr)
	}

	printTableFooter()
	fmt.Println()

	// Footer with default model
	if summary.DefaultModelID != "" {
		for _, model := range models {
			if model.ModelID == summary.DefaultModelID {
				color.Cyan("  Default Model: %s", model.DisplayName)
				break
			}
		}
		fmt.Println()
	}
}

// printTableHeader prints the table header
func printTableHeader() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(cyan("  ┌────────────────────────┬──────────────────────┬─────────────┬──────────┐"))
	fmt.Printf("%s%-24s%s%-22s%s%-13s%s%-10s%s\n",
		cyan("  │ "),
		cyan("Model"),
		cyan(" │ "),
		cyan("Quota"),
		cyan(" │ "),
		cyan("Reset In"),
		cyan(" │ "),
		cyan("Status"),
		cyan(" │"))
	fmt.Println(cyan("  ├────────────────────────┼──────────────────────┼─────────────┼──────────┤"))
}

// printTableRow prints a table row
func printTableRow(model, quota, reset, status string) {
	// Remove ANSI color codes for width calculation
	statusClean := strings.ReplaceAll(status, "\x1b[32m", "")
	statusClean = strings.ReplaceAll(statusClean, "\x1b[33m", "")
	statusClean = strings.ReplaceAll(statusClean, "\x1b[31m", "")
	statusClean = strings.ReplaceAll(statusClean, "\x1b[0m", "")

	// Calculate padding for status (accounting for unicode characters)
	statusWidth := 6 // Width of "✓ OK", "⚠ LOW", "✗ EMPTY"
	statusPadding := strings.Repeat(" ", 10-statusWidth)

	fmt.Printf("  │ %-22s │ %-20s │ %-11s │ %s%s │\n",
		truncate(model, 22),
		quota,
		reset,
		status,
		statusPadding,
	)
}

// printTableFooter prints the table footer
func printTableFooter() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Println(cyan("  └────────────────────────┴──────────────────────┴─────────────┴──────────┘"))
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// formatQuota formats the quota percentage with a progress bar
func formatQuota(model models.ModelQuota) string {
	percentage := model.GetRemainingPercentage()

	// Create progress bar (10 characters)
	filled := percentage / 10
	empty := 10 - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	return fmt.Sprintf("%3d%% [%s]", percentage, bar)
}

// formatResetTime formats the time until reset in a human-readable format
func formatResetTime(model models.ModelQuota) string {
	duration := model.GetTimeUntilReset()

	if duration < 0 {
		return "Resetting..."
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	return fmt.Sprintf("%dm", minutes)
}

// formatStatus formats the status with appropriate coloring
func formatStatus(model models.ModelQuota) string {
	status := model.GetStatusString()

	switch status {
	case "OK":
		return color.GreenString("✓ OK")
	case "LOW":
		return color.YellowString("⚠ LOW")
	case "EMPTY":
		return color.RedString("✗ EMPTY")
	default:
		return status
	}
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
