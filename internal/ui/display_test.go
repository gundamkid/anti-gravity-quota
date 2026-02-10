package ui

import (
	"testing"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

func TestFormatResetTime(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		reset    time.Time
		expected string
	}{
		{
			name:     "Regenerating",
			reset:    now.Add(-time.Hour),
			expected: "Regenerating...",
		},
		{
			name:     "Minutes",
			reset:    now.Add(5 * time.Minute),
			expected: "5m",
		},
		{
			name:     "Hours and minutes",
			reset:    now.Add(2*time.Hour + 30*time.Minute),
			expected: "2h 30m",
		},
		{
			name:     "Days and hours",
			reset:    now.Add(48*time.Hour + 2*time.Hour),
			expected: "2d 02h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := models.ModelQuota{ResetTime: tt.reset}
			got := formatResetTime(m, now)
			if got != tt.expected {
				t.Errorf("formatResetTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSpinner(t *testing.T) {
	s := NewSpinner()
	f1 := s.Next()
	f2 := s.Next()
	if f1 == f2 {
		t.Error("spinner did not advance")
	}
	if f1 == "" || f2 == "" {
		t.Error("spinner frame is empty")
	}
}

func TestShortenModelName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Claude Opus 3.5 (Thinking)", "Opus 3.5 (T)"},
		{"Claude Sonnet 3.5", "Sonnet 3.5"},
		{"Gemini 1.5 Flash", "Gem 1.5 Flash"},
		{"Model (Thinking)", "Model (T)"},
		{"Model (Low)", "Model (L)"},
		{"Model (Medium)", "Model (M)"},
		{"Model (High)", "Model (H)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shortenModelName(tt.input)
			if got != tt.expected {
				t.Errorf("shortenModelName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
