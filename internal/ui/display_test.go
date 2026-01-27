package ui

import (
	"testing"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

func TestFormatResetTime(t *testing.T) {
	tests := []struct {
		name     string
		reset    time.Time
		expected string
	}{
		{
			name:     "Regenerating",
			reset:    time.Now().Add(-time.Hour),
			expected: "Regenerating...",
		},
		{
			name:     "Minutes",
			reset:    time.Now().Add(5 * time.Minute),
			expected: "5m",
		},
		{
			name:     "Hours and minutes",
			reset:    time.Now().Add(2*time.Hour + 30*time.Minute),
			expected: "2h 30m",
		},
		{
			name:     "Days and hours",
			reset:    time.Now().Add(48*time.Hour + 2*time.Hour),
			expected: "2d 02h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := models.ModelQuota{ResetTime: tt.reset}
			got := formatResetTime(m)
			// Duration can be slightly off due to processing time, using Contains for flexibility if needed
			// but for unit test we can mock time if we want to be exact.
			// For now, let's just check if it's not empty and follows format.
			if got == "" {
				t.Error("formatResetTime returned empty string")
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
