package notify

import (
	"testing"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

func TestStateTracker(t *testing.T) {
	tracker := NewStateTracker()
	email := "test@example.com"
	resetTime := time.Now().Add(1 * time.Hour)

	t.Run("First Fetch - Baseline Summary", func(t *testing.T) {
		quotas := []models.ModelQuota{
			{DisplayName: "claude-3-opus", RemainingFraction: 1.0},                        // HEALTHY
			{DisplayName: "gemini-1.5-pro", RemainingFraction: 0.1, ResetTime: resetTime}, // CRITICAL
		}

		changes := tracker.Update(email, quotas)

		// Should notify ALL models on first fetch for initial summary
		if len(changes) != 2 {
			t.Errorf("expected 2 changes, got %d", len(changes))
		}

		// Check statuses
		if changes[0].OldStatus != "INITIAL" || changes[1].OldStatus != "INITIAL" {
			t.Error("expected INITIAL old status on first fetch")
		}

		// Second change should be the gemini one (order preserved)
		if changes[1].DisplayName != "gemini-1.5-pro" || changes[1].NewStatus != "CRITICAL" {
			t.Errorf("unexpected second change: %+v", changes[1])
		}
	})

	t.Run("Second Fetch - No Change", func(t *testing.T) {
		quotas := []models.ModelQuota{
			{DisplayName: "claude-3-opus", RemainingFraction: 1.0},
			{DisplayName: "gemini-1.5-pro", RemainingFraction: 0.1, ResetTime: resetTime},
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
				if c.OldPercentage != 100 || c.NewPercentage != 40 {
					t.Errorf("unexpected percentages: old=%d, new=%d", c.OldPercentage, c.NewPercentage)
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
		// Should treat as first fetch again, returning the model
		if len(changes) != 1 {
			t.Errorf("expected 1 baseline notification after reset, got %d", len(changes))
		}
		if changes[0].OldStatus != "INITIAL" {
			t.Errorf("expected INITIAL status after reset, got %s", changes[0].OldStatus)
		}
	})

	t.Run("Skip Empty Names", func(t *testing.T) {
		tracker.Reset()
		quotas := []models.ModelQuota{
			{DisplayName: "", RemainingFraction: 0.1},
			{DisplayName: "valid-model", RemainingFraction: 0.5},
		}

		changes := tracker.Update(email, quotas)

		if len(changes) != 1 {
			t.Errorf("expected 1 change, got %d", len(changes))
		}
		if changes[0].DisplayName != "valid-model" {
			t.Errorf("expected valid-model, got %s", changes[0].DisplayName)
		}
	})
}
