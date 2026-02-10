package models

import (
	"testing"
	"time"
)

func TestModelQuota_GetStatusString(t *testing.T) {
	tests := []struct {
		name     string
		quota    ModelQuota
		expected string
	}{
		{
			name: "Exhausted (Flag)",
			quota: ModelQuota{
				IsExhausted: true,
			},
			expected: "EMPTY",
		},
		{
			name: "Exhausted (Fraction)",
			quota: ModelQuota{
				RemainingFraction: 0,
			},
			expected: "EMPTY",
		},
		{
			name: "Critical quota (10%)",
			quota: ModelQuota{
				RemainingFraction: 0.1,
			},
			expected: "CRITICAL",
		},
		{
			name: "Critical quota (20%)",
			quota: ModelQuota{
				RemainingFraction: 0.2,
			},
			expected: "CRITICAL",
		},
		{
			name: "Warning quota (25%)",
			quota: ModelQuota{
				RemainingFraction: 0.25,
			},
			expected: "WARNING",
		},
		{
			name: "Warning quota (50%)",
			quota: ModelQuota{
				RemainingFraction: 0.5,
			},
			expected: "WARNING",
		},
		{
			name: "Healthy quota (60%)",
			quota: ModelQuota{
				RemainingFraction: 0.6,
			},
			expected: "HEALTHY",
		},
		{
			name: "Healthy quota (100%)",
			quota: ModelQuota{
				RemainingFraction: 1.0,
			},
			expected: "HEALTHY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.quota.GetStatusString() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.quota.GetStatusString())
			}
		})
	}
}

func TestModelQuota_GetRemainingPercentage(t *testing.T) {
	quota := ModelQuota{RemainingFraction: 0.75}
	if quota.GetRemainingPercentage() != 75 {
		t.Errorf("expected 75, got %d", quota.GetRemainingPercentage())
	}
}

func TestGetDefaultMetadata(t *testing.T) {
	meta := GetDefaultMetadata()
	if meta.IDEType != "ANTIGRAVITY" {
		t.Errorf("expected ANTIGRAVITY, got %s", meta.IDEType)
	}
}

func TestModel_ToModelQuota(t *testing.T) {
	resetTime := time.Now().Add(time.Hour)
	m := Model{
		DisplayName: "Test Model",
		Label:       "test-label",
		QuotaInfo: ModelQuotaInfo{
			RemainingFraction: 0.8,
			ResetTime:         resetTime,
			IsExhausted:       false,
		},
		ModelProvider: "GOOGLE",
	}

	mq := m.ToModelQuota("model-1")

	if mq.ModelID != "model-1" || mq.DisplayName != "Test Model" || mq.RemainingFraction != 0.8 {
		t.Errorf("ModelQuota conversion failed: %+v", mq)
	}
	if !mq.ResetTime.Equal(resetTime) {
		t.Errorf("ResetTime mismatch")
	}
}
