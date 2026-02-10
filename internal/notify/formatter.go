package notify

import (
	"fmt"
	"strings"
	"time"
)

// MessageFormatter handles building notification messages
type MessageFormatter struct{}

// NewMessageFormatter creates a new message formatter
func NewMessageFormatter() *MessageFormatter {
	return &MessageFormatter{}
}

// FormatChanges aggregates multiple status changes into a single notification message grouped by Account.
func (f *MessageFormatter) FormatChanges(changes []StatusChange) Message {
	if len(changes) == 0 {
		return Message{}
	}

	// Determine overall severity
	maxSeverity := SeverityInfo
	isInitial := false
	for _, c := range changes {
		severity := f.getSeverity(c.NewStatus)
		if severity > maxSeverity {
			maxSeverity = severity
		}
		if c.OldStatus == "INITIAL" {
			isInitial = true
		}
	}

	title := "🔄 Status Update"
	if isInitial {
		title = "📊 Quota Summary"
	}

	// Group by Account -> Status
	type Grouped struct {
		Account  string
		ByStatus map[string][]StatusChange
	}

	var accounts []string
	accountGroups := make(map[string]map[string][]StatusChange)

	for _, c := range changes {
		acc := c.Account
		if acc == "" {
			acc = "Unknown Account"
		}
		if _, exists := accountGroups[acc]; !exists {
			accounts = append(accounts, acc)
			accountGroups[acc] = make(map[string][]StatusChange)
		}
		accountGroups[acc][c.NewStatus] = append(accountGroups[acc][c.NewStatus], c)
	}

	var sb strings.Builder
	statusOrder := []string{"HEALTHY", "WARNING", "CRITICAL", "EMPTY"}

	for i, acc := range accounts {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("👤 *%s*\n", acc))

		group := accountGroups[acc]
		for _, status := range statusOrder {
			items := group[status]
			if len(items) == 0 {
				continue
			}

			sb.WriteString(fmt.Sprintf("  %s\n", f.getStatusHeader(status)))
			for _, c := range items {
				// Base line: - Model Name | X%
				line := fmt.Sprintf("    - %s | %d%%", c.DisplayName, c.NewPercentage)

				// Delta logic for updates: (↓ 70%)
				if !isInitial && c.OldStatus != "UNKNOWN" && c.OldStatus != "INITIAL" {
					delta := c.NewPercentage - c.OldPercentage
					if delta < 0 {
						line += fmt.Sprintf(" (↓ %d%%)", -delta)
					} else if delta > 0 {
						line += fmt.Sprintf(" (↑ %d%%)", delta)
					}
				}

				// Reset time for Critical/Empty: ⏳ 2h 30m
				if (status == "EMPTY" || status == "CRITICAL") && !c.ResetTime.IsZero() {
					remaining := time.Until(c.ResetTime)
					if remaining > 0 {
						line += fmt.Sprintf(" ⏳ %s", FormatTimeRemaining(remaining))
					}
				}

				sb.WriteString(line + "\n")
			}
		}
	}

	return Message{
		Title:    title,
		Body:     strings.TrimSpace(sb.String()),
		Severity: maxSeverity,
	}
}

func (f *MessageFormatter) getSeverity(status string) Severity {
	switch status {
	case "EMPTY":
		return SeverityCritical
	case "CRITICAL":
		return SeverityCritical
	case "WARNING":
		return SeverityWarning
	case "HEALTHY":
		return SeverityRecovery
	default:
		return SeverityInfo
	}
}

func (f *MessageFormatter) getStatusHeader(status string) string {
	switch status {
	case "EMPTY":
		return "❌ *Empty*"
	case "CRITICAL":
		return "⛔ *Critical*"
	case "WARNING":
		return "⚠️ *Warning*"
	case "HEALTHY":
		return "✅ *Healthy*"
	default:
		return "*" + status + "*"
	}
}

// FormatTimeRemaining returns a human readable time duration
func FormatTimeRemaining(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
