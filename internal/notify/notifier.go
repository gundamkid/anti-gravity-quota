package notify

import (
	"context"
	"sync"
)

// Severity represents the importance of the notification
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
	SeverityRecovery Severity = "RECOVERY"
)

// Message represents a notification message
type Message struct {
	Title    string
	Body     string
	Severity Severity
}

// Notifier is the interface that all notification channels must implement
type Notifier interface {
	// Name returns the identifier of the notifier (e.g., "telegram")
	Name() string
	// Send sends a notification message
	Send(ctx context.Context, msg Message) error
	// IsEnabled returns whether the notifier is configured and enabled
	IsEnabled() bool
}

// Registry manages multiple notifiers
type Registry struct {
	mu        sync.RWMutex
	notifiers map[string]Notifier
}

// NewRegistry creates a new notifier registry
func NewRegistry() *Registry {
	return &Registry{
		notifiers: make(map[string]Notifier),
	}
}

// Register adds a notifier to the registry
func (r *Registry) Register(n Notifier) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifiers[n.Name()] = n
}

// Get returns a notifier by name
func (r *Registry) Get(name string) (Notifier, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.notifiers[name]
	return n, ok
}

// NotifyAll sends a message to all enabled notifiers
func (r *Registry) NotifyAll(ctx context.Context, msg Message) []error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var errs []error
	for _, n := range r.notifiers {
		if n.IsEnabled() {
			if err := n.Send(ctx, msg); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// List returns names of all registered notifiers
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var names []string
	for name := range r.notifiers {
		names = append(names, name)
	}
	return names
}
