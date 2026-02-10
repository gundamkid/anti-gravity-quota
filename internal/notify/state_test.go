package notify

import (
	"testing"

	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

func TestStateTracker(t *testing.T) {
	tracker := NewStateTracker()
	email := "test@example.com"

	t.Run("First Fetch - Baseline", func(t *testing.T) {
		quotas := []models.ModelQuota{
			{DisplayName: "claude-3-opus", RemainingFraction: 1.0},  // HEALTHY
			{DisplayName: "gemini-1.5-pro", RemainingFraction: 0.1}, // CRITICAL
		}

		changes := tracker.Update(email, quotas)

		// Should only notify the non-healthy model
		if len(changes) != 1 {
			t.Errorf("expected 1 change, got %d", len(changes))
		}
		if changes[0].DisplayName != "gemini-1.5-pro" || changes[0].NewStatus != "CRITICAL" {
			t.Errorf("unexpected change: %+v", changes[0])
		}
	})

	t.Run("Second Fetch - No Change", func(t *testing.T) {
		quotas := []models.ModelQuota{
			{DisplayName: "claude-3-opus", RemainingFraction: 1.0},
			{DisplayName: "gemini-1.5-pro", RemainingFraction: 0.1},
		}

		changes := tracker.Update(email, quotas)

		if len(changes) != 0 {
			t.Errorf("expected 0 changes, got %d", len(changes))
		}
	})

	t.Run("Third Fetch - Status Changed", func(t *testing.T) {
		quotas := []models.ModelQuota{
			{DisplayName: "claude-3-opus", RemainingFraction: 0.4},  // WARNING
			{DisplayName: "gemini-1.5-pro", RemainingFraction: 0.8}, // HEALTHY (Recovery)
		}

		changes := tracker.Update(email, quotas)

		if len(changes) != 2 {
			t.Errorf("expected 2 changes, got %d", len(changes))
		}

		// Check Claude
		var claudeFound bool
		for _, c := range changes {
			if c.DisplayName == "claude-3-opus" {
				claudeFound = true
				if c.OldStatus != "HEALTHY" || c.NewStatus != "WARNING" {
					t.Errorf("unexpected claude change: %+v", c)
				}
			}
		}
		if !claudeFound {
			t.Error("claude change not found")
		}
	})

	t.Run("Reset", func(t *testing.T) {
		tracker.Reset()
		quotas := []models.ModelQuota{
			{DisplayName: "claude-3-opus", RemainingFraction: 0.4},
		}
		changes := tracker.Update(email, quotas)
		// Should treat as first fetch again
		if len(changes) != 1 {
			t.Error("expected baseline notification after reset")
		}
	})
}
