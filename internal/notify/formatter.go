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

// FormatChanges aggregates multiple status changes into a single notification message
func (f *MessageFormatter) FormatChanges(changes []StatusChange) Message {
	if len(changes) == 0 {
		return Message{}
	}

	// Determine overall severity based on the worst change
	maxSeverity := SeverityInfo
	for _, c := range changes {
		severity := f.getSeverity(c.NewStatus)
		if severity > maxSeverity {
			maxSeverity = severity
		}
	}

	title := "[AG-Quota] Status Update"
	if len(changes) == 1 {
		title = fmt.Sprintf("[AG-Quota] %s: %s", changes[0].DisplayName, changes[0].NewStatus)
	}

	var sb strings.Builder

	// Group by status
	byStatus := make(map[string][]StatusChange)
	for _, c := range changes {
		byStatus[c.NewStatus] = append(byStatus[c.NewStatus], c)
	}

	// Order of display
	order := []string{"EMPTY", "CRITICAL", "WARNING", "HEALTHY"}

	for _, status := range order {
		cmds := byStatus[status]
		if len(cmds) == 0 {
			continue
		}

		sb.WriteString(f.getStatusHeader(status) + "\n")
		for _, c := range cmds {
			sb.WriteString(fmt.Sprintf("• %s: %s (from %s)\n", c.DisplayName, c.NewStatus, c.OldStatus))
			if c.Account != "" {
				sb.WriteString(fmt.Sprintf("  └ Account: %s\n", c.Account))
			}
		}
		sb.WriteString("\n")
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
		return "🚫 *EMPTY*"
	case "CRITICAL":
		return "🔴 *CRITICAL*"
	case "WARNING":
		return "⚠️ *WARNING*"
	case "HEALTHY":
		return "✅ *RECOVERED*"
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
