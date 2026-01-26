package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

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
func DisplayQuotaSummary(summary *models.QuotaSummary) {
	// Header
	fmt.Println()
	color.Cyan("  ✨ Anti-Gravity Quota Status")
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

	t.AppendHeader(table.Row{"Model", "Quota", "Reset In", "Status"})

	for _, model := range models {
		percentage := model.GetRemainingPercentage()

		// Colorize Quota cell
		var quotaColor text.Colors
		if percentage <= 10 {
			quotaColor = text.Colors{text.FgRed, text.Bold}
		} else if percentage <= 30 {
			quotaColor = text.Colors{text.FgYellow}
		} else {
			quotaColor = text.Colors{text.FgGreen}
		}
		quotaStr := fmt.Sprintf("%3d%%", percentage)

		// Format Status with colors
		statusStr := model.GetStatusString()
		var statusColor text.Colors
		switch statusStr {
		case "OK":
			statusStr = "✓ OK"
			statusColor = text.Colors{text.FgGreen}
		case "LOW":
			statusStr = "⚠ LOW"
			statusColor = text.Colors{text.FgYellow}
		case "EMPTY":
			statusStr = "✗ EMPTY"
			statusColor = text.Colors{text.FgRed}
		}

		t.AppendRow(table.Row{
			model.DisplayName,
			quotaColor.Sprint(quotaStr),
			formatResetTime(model),
			statusColor.Sprint(statusStr),
		})
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
				color.Cyan("  ⭐ Default Model: %s", model.DisplayName)
				break
			}
		}
		fmt.Println()
	}
}

// formatResetTime formats the time until reset in a human-readable format
func formatResetTime(model models.ModelQuota) string {
	duration := model.GetTimeUntilReset()

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
