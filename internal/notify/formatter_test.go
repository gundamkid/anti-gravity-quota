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

		if msg.Title != "🔄 Status Update" {
			t.Errorf("wrong title, got %s", msg.Title)
		}
		if !strings.Contains(msg.Body, "⚠️ *Warning*") {
			t.Error("body should contain warning header")
		}
		// Check delta formatting: 45% (↓ 55%)
		if !strings.Contains(msg.Body, "45% (↓ 55%)") {
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

		if !strings.Contains(msg.Body, "⛔ *Critical*") {
			t.Error("body should contain critical header")
		}
		if !strings.Contains(msg.Body, "⏳ 2h 30m") {
			t.Errorf("body missing reset time, got: %s", msg.Body)
		}
	})

	t.Run("Multiple Changes - Batching and Grouping", func(t *testing.T) {
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
				Account:       "acc1@gmail.com",
				DisplayName:   "Model B",
				OldStatus:     "WARNING",
				NewStatus:     "HEALTHY",
				OldPercentage: 20,
				NewPercentage: 100,
			},
		}

		msg := formatter.FormatChanges(changes)

		if !strings.Contains(msg.Body, "👤 *acc1@gmail.com*") {
			t.Error("body should contain account header")
		}
		if !strings.Contains(msg.Body, "⛔ *Critical*") {
			t.Error("body should contain critical header")
		}
		if !strings.Contains(msg.Body, "✅ *Healthy*") {
			t.Error("body should contain healthy header")
		}
		// Verify Healthy comes AFTER Critical in update order if both present?
		// Actually statusOrder is HEALTHY, WARNING, CRITICAL, EMPTY.
		// So Healthy should be first.
		healthyIdx := strings.Index(msg.Body, "✅ *Healthy*")
		criticalIdx := strings.Index(msg.Body, "⛔ *Critical*")
		if healthyIdx > criticalIdx {
			t.Error("Healthy should come before Critical in our defined order")
		}
	})

	t.Run("Initial Summary", func(t *testing.T) {
		changes := []StatusChange{
			{
				Account:       "user@gmail.com",
				DisplayName:   "Gemini 1.5 Flash",
				OldStatus:     "INITIAL",
				NewStatus:     "HEALTHY",
				NewPercentage: 100,
			},
			{
				Account:       "user@gmail.com",
				DisplayName:   "Claude 3.5 Sonnet",
				OldStatus:     "INITIAL",
				NewStatus:     "CRITICAL",
				NewPercentage: 5,
			},
		}

		msg := formatter.FormatChanges(changes)

		if msg.Title != "📊 Quota Summary" {
			t.Errorf("wrong title for initial summary: %s", msg.Title)
		}
		if strings.Contains(msg.Body, "INITIAL") {
			t.Error("body should not contain INITIAL sentinel")
		}
		if strings.Contains(msg.Body, "↓") || strings.Contains(msg.Body, "↑") {
			t.Error("body should not show deltas for initial summary")
		}
		if !strings.Contains(msg.Body, "⛔ *Critical*") || !strings.Contains(msg.Body, "✅ *Healthy*") {
			t.Error("body missing headers")
		}
	})

	t.Run("Alphabetical Sorting", func(t *testing.T) {
		changes := []StatusChange{
			{
				Account:       "user@gmail.com",
				DisplayName:   "Z-Model",
				NewStatus:     "HEALTHY",
				NewPercentage: 100,
			},
			{
				Account:       "user@gmail.com",
				DisplayName:   "A-Model",
				NewStatus:     "HEALTHY",
				NewPercentage: 100,
			},
			{
				Account:       "user@gmail.com",
				DisplayName:   "M-Model",
				NewStatus:     "HEALTHY",
				NewPercentage: 100,
			},
		}

		msg := formatter.FormatChanges(changes)

		aIdx := strings.Index(msg.Body, "A-Model")
		mIdx := strings.Index(msg.Body, "M-Model")
		zIdx := strings.Index(msg.Body, "Z-Model")

		if aIdx == -1 || mIdx == -1 || zIdx == -1 {
			t.Fatal("models missing from body")
		}

		if !(aIdx < mIdx && mIdx < zIdx) {
			t.Errorf("models not sorted alphabetically: A-Model(%d), M-Model(%d), Z-Model(%d)", aIdx, mIdx, zIdx)
		}
	})
}
