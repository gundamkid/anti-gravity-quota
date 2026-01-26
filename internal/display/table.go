package display

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gundamkid/anti-gravity-quota/internal/client"
	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"golang.org/x/term"
)

var (
	// Colors
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F5587B"}
	text      = lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}

	// Borders
	border = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	// Styles
	boxStyle = lipgloss.NewStyle().
			Border(border).
			BorderForeground(highlight).
			Padding(1, 0)

	headerStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			Align(lipgloss.Center)

	connInfoStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginBottom(1)
)

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%dh %dm", h, m)
}

func progressBar(width int, percent float64) string {
	w := float64(width)
	fullSize := int(math.Round(w * percent / 100))
	var fullCells string
	
	// Gradient effect for progress bar based on percentage
	var barColor lipgloss.Color
	if percent < 50 {
		barColor = lipgloss.Color("#43BF6D") // Green
	} else if percent < 80 {
		barColor = lipgloss.Color("#F0AF3A") // Yellow
	} else {
		barColor = lipgloss.Color("#F5587B") // Red
	}

	fullCells = lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", fullSize))
	emptyCells := strings.Repeat("░", int(w)-fullSize)
	
	return fmt.Sprintf("%s%s", fullCells, emptyCells)
}

func RenderTable(status *client.UserStatusResponse, cfg *config.Config) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if width == 0 {
		width = 80
	}
	// Cap max width for readability
	if width > 100 {
		width = 100
	}

	accountEmail := "Unknown Account"
	active := cfg.GetActiveAccount()
	if active != nil {
		accountEmail = active.Email
	}

	doc := strings.Builder{}

	// Header
	header := headerStyle.Width(width - 4).Render("Anti-Gravity Quota Monitor")
	doc.WriteString(header + "\n")
	doc.WriteString(connInfoStyle.Width(width - 4).Align(lipgloss.Center).Render(fmt.Sprintf("Account: %s", accountEmail)) + "\n\n")

	// Table Header
	// Model | Used/Limit | Remaining Bar | Resets In
	
	// Columns with fixed ratios equivalent
	colModelWidth := 20
	colUsedWidth := 18
	colResetWidth := 15
	colBarWidth := width - 4 - colModelWidth - colUsedWidth - colResetWidth - 6 // padding/borders
	if colBarWidth < 10 {
		colBarWidth = 10
	}

	rowStyle := lipgloss.NewStyle().Padding(0, 1)

	headerRow := rowStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Width(colModelWidth).Foreground(subtle).Render("MODEL"),
			lipgloss.NewStyle().Width(colUsedWidth).Foreground(subtle).Render("USED/LIMIT"),
			lipgloss.NewStyle().Width(colBarWidth).Foreground(subtle).Render("REMAINING"),
			lipgloss.NewStyle().Width(colResetWidth).Foreground(subtle).Render("RESETS IN"),
		),
	)
	doc.WriteString(headerRow + "\n")
	doc.WriteString(lipgloss.NewStyle().Foreground(subtle).Render(strings.Repeat("─", width-2)) + "\n")

	for _, model := range status.Models {
		percentUsed := 0.0
		if model.QuotaLimit > 0 {
			percentUsed = (float64(model.QuotaUsed) / float64(model.QuotaLimit)) * 100
		}
		
		percentRemaining := 100.0 - percentUsed
		if percentRemaining < 0 { percentRemaining = 0 }

		resetIn := time.Until(model.ResetAt)
		if resetIn < 0 {
			resetIn = 0
		}

		row := rowStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Width(colModelWidth).Foreground(highlight).Render(model.Name),
				lipgloss.NewStyle().Width(colUsedWidth).Render(fmt.Sprintf("%d/%d", model.QuotaUsed, model.QuotaLimit)),
				lipgloss.NewStyle().Width(colBarWidth).Render(fmt.Sprintf("%3.0f%% %s", percentRemaining, progressBar(colBarWidth-6, percentRemaining/100*100))), // Note: progress bar here visuals 'remaining' or 'used'? Usually progress bar shows used?
				// User guide says: Gemini 3 Pro | 1,500/5,000 | 70% [====...] | 4h 32m
				// Let's implement showing REMAINING percentage and bar represents REMAINING capacity?
				// Or USED?
				// Guide: "70% [XXX...]"
				// If 1500 used of 5000, that's 30% used, 70% remaining.
				// So the bar should probably be full if fresh, and empty if used up.
				// Let's make bar represent remaining (Green full -> Red empty).
				lipgloss.NewStyle().Width(colResetWidth).Render(formatDuration(resetIn)),
			),
		)
		doc.WriteString(row + "\n")
	}

	fmt.Println(boxStyle.Width(width).Render(doc.String()))
}

// RenderCompact renders a single line summary
func RenderCompact(status *client.UserStatusResponse) {
	var parts []string
	for _, model := range status.Models {
		percentUsed := 0.0
		if model.QuotaLimit > 0 {
			percentUsed = (float64(model.QuotaUsed) / float64(model.QuotaLimit)) * 100
		}
		remaining := 100.0 - percentUsed
		
		// Colorize based on remaining
		style := lipgloss.NewStyle()
		if remaining < 20 {
			style = style.Foreground(warning)
		} else {
			style = style.Foreground(special)
		}
		
		// Simplify name: "Gemini 1.5 Pro" -> "Gemini"
		name := strings.Split(model.Name, " ")[0]
		parts = append(parts, fmt.Sprintf("%s: %s", name, style.Render(fmt.Sprintf("%.0f%%", remaining))))
	}
	
	// Add reset time of first model (assuming synced)
	if len(status.Models) > 0 {
		resetIn := time.Until(status.Models[0].ResetAt)
		parts = append(parts, fmt.Sprintf("Reset: %s", formatDuration(resetIn)))
	}
	
	fmt.Println(strings.Join(parts, " | "))
}

// RenderJSON simply marshals to JSON
func RenderJSON(status *client.UserStatusResponse) {
    // Handled in main usually, but helper here for consistency?
    // Nah, main can handle json marshalling trivially.
}
