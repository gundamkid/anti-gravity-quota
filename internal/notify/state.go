package notify

import (
	"sync"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

// StatusChange represents a change in quota status for a model
type StatusChange struct {
	Account       string
	DisplayName   string
	OldStatus     string
	NewStatus     string
	OldPercentage int
	NewPercentage int
	ResetTime     time.Time
}

// StateTracker monitors status changes between fetches
type StateTracker struct {
	mu sync.RWMutex
	// lastStatus stores [accountEmail][displayName] = statusString
	lastStatus map[string]map[string]string
	// lastPercentage stores [accountEmail][displayName] = percentage
	lastPercentage map[string]map[string]int
	// isFirstFetch tracks if we have baseline data for an account
	isFirstFetch map[string]bool
}

// NewStateTracker creates a new status state tracker
func NewStateTracker() *StateTracker {
	return &StateTracker{
		lastStatus:     make(map[string]map[string]string),
		lastPercentage: make(map[string]map[string]int),
		isFirstFetch:   make(map[string]bool),
	}
}

// Update updates the state for an account and returns detected status changes.
// If it's the first time seeing this account, it returns all non-HEALTHY statuses as changes.
func (t *StateTracker) Update(accountEmail string, quotas []models.ModelQuota) []StatusChange {
	t.mu.Lock()
	defer t.mu.Unlock()

	var changes []StatusChange

	if t.lastStatus[accountEmail] == nil {
		t.lastStatus[accountEmail] = make(map[string]string)
		t.lastPercentage[accountEmail] = make(map[string]int)
		t.isFirstFetch[accountEmail] = true
	}

	isFirst := t.isFirstFetch[accountEmail]
	t.isFirstFetch[accountEmail] = false

	for _, q := range quotas {
		displayName := q.DisplayName
		newStatus := q.GetStatusString()
		newPercentage := q.GetRemainingPercentage()
		oldStatus, exists := t.lastStatus[accountEmail][displayName]
		oldPercentage := t.lastPercentage[accountEmail][displayName]

		if isFirst {
			// On first fetch, notify if not HEALTHY
			if newStatus != "HEALTHY" {
				changes = append(changes, StatusChange{
					Account:       accountEmail,
					DisplayName:   displayName,
					OldStatus:     "UNKNOWN",
					NewStatus:     newStatus,
					NewPercentage: newPercentage,
					ResetTime:     q.ResetTime,
				})
			}
		} else if exists && oldStatus != newStatus {
			// Status changed
			changes = append(changes, StatusChange{
				Account:       accountEmail,
				DisplayName:   displayName,
				OldStatus:     oldStatus,
				NewStatus:     newStatus,
				OldPercentage: oldPercentage,
				NewPercentage: newPercentage,
				ResetTime:     q.ResetTime,
			})
		}

		// Update state
		t.lastStatus[accountEmail][displayName] = newStatus
		t.lastPercentage[accountEmail][displayName] = newPercentage
	}

	return changes
}

// Reset clears the state tracker
func (t *StateTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastStatus = make(map[string]map[string]string)
	t.lastPercentage = make(map[string]map[string]int)
	t.isFirstFetch = make(map[string]bool)
}
