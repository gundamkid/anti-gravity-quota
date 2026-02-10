package notify

import (
	"strings"
	"testing"
	"time"
)

func TestMessageFormatter(t *testing.T) {
	formatter := NewMessageFormatter()

	t.Run("Single Change with Delta", func(t *testing.T) {
		changes := []StatusChange{
			{
				Account:       "user@gmail.com",
				DisplayName:   "Claude 3 Opus",
				OldStatus:     "HEALTHY",
				NewStatus:     "WARNING",
				OldPercentage: 100,
				NewPercentage: 45,
			},
		}

		msg := formatter.FormatChanges(changes)

		if !strings.Contains(msg.Title, "Claude 3 Opus") {
			t.Errorf("title should contain model name, got %s", msg.Title)
		}
		if !strings.Contains(msg.Body, "⚠️ *WARNING*") {
			t.Error("body should contain warning header")
		}
		// Check delta formatting: (100% ↓ 45% (↓55%))
		if !strings.Contains(msg.Body, "100% ↓ 45% (↓55%)") {
			t.Errorf("body missing delta, got: %s", msg.Body)
		}
		if msg.Severity != SeverityWarning {
			t.Errorf("expected warning severity, got %v", msg.Severity)
		}
	})

	t.Run("Critical with Reset Time", func(t *testing.T) {
		resetTime := time.Now().Add(2*time.Hour + 30*time.Minute)
		changes := []StatusChange{
			{
				Account:       "acc1@gmail.com",
				DisplayName:   "Model A",
				OldStatus:     "HEALTHY",
				NewStatus:     "CRITICAL",
				OldPercentage: 80,
				NewPercentage: 5,
				ResetTime:     resetTime,
			},
		}

		msg := formatter.FormatChanges(changes)

		if !strings.Contains(msg.Body, "🔴 *CRITICAL*") {
			t.Error("body should contain critical header")
		}
		if !strings.Contains(msg.Body, "Reset in 2h 30m") {
			t.Errorf("body missing reset time, got: %s", msg.Body)
		}
		if msg.Severity != SeverityCritical {
			t.Errorf("overall severity should be critical, got %v", msg.Severity)
		}
	})

	t.Run("Multiple Changes - Batching", func(t *testing.T) {
		changes := []StatusChange{
			{
				Account:       "acc1@gmail.com",
				DisplayName:   "Model A",
				OldStatus:     "HEALTHY",
				NewStatus:     "CRITICAL",
				OldPercentage: 80,
				NewPercentage: 5,
			},
			{
				Account:       "acc2@gmail.com",
				DisplayName:   "Model B",
				OldStatus:     "WARNING",
				NewStatus:     "HEALTHY",
				OldPercentage: 20,
				NewPercentage: 100,
			},
		}

		msg := formatter.FormatChanges(changes)

		if !strings.Contains(msg.Body, "🔴 *CRITICAL*") {
			t.Error("body should contain critical header")
		}
		if !strings.Contains(msg.Body, "✅ *RECOVERED*") {
			t.Error("body should contain recovered header")
		}
		if !strings.Contains(msg.Body, "Account: acc1@gmail.com") {
			t.Error("body should contain account email")
		}
	})
}
