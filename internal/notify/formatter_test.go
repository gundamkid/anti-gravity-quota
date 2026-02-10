package notify

import (
	"strings"
	"testing"
)

func TestMessageFormatter(t *testing.T) {
	formatter := NewMessageFormatter()

	t.Run("Single Change", func(t *testing.T) {
		changes := []StatusChange{
			{
				Account:     "user@gmail.com",
				DisplayName: "Claude 3 Opus",
				OldStatus:   "HEALTHY",
				NewStatus:   "WARNING",
			},
		}

		msg := formatter.FormatChanges(changes)

		if !strings.Contains(msg.Title, "Claude 3 Opus") {
			t.Errorf("title should contain model name, got %s", msg.Title)
		}
		if !strings.Contains(msg.Body, "⚠️ *WARNING*") {
			t.Error("body should contain warning header")
		}
		if msg.Severity != SeverityWarning {
			t.Errorf("expected warning severity, got %v", msg.Severity)
		}
	})

	t.Run("Multiple Changes - Batching", func(t *testing.T) {
		changes := []StatusChange{
			{
				Account:     "acc1@gmail.com",
				DisplayName: "Model A",
				OldStatus:   "HEALTHY",
				NewStatus:   "CRITICAL",
			},
			{
				Account:     "acc2@gmail.com",
				DisplayName: "Model B",
				OldStatus:   "WARNING",
				NewStatus:   "HEALTHY",
			},
		}

		msg := formatter.FormatChanges(changes)

		if !strings.Contains(msg.Body, "🔴 *CRITICAL*") {
			t.Error("body should contain critical header")
		}
		if !strings.Contains(msg.Body, "✅ *RECOVERED*") {
			t.Error("body should contain recovered header")
		}
		if msg.Severity != SeverityCritical {
			t.Errorf("overall severity should be critical, got %v", msg.Severity)
		}
	})
}
