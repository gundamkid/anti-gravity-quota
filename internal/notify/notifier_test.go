package notify

import (
	"context"
	"errors"
	"testing"
)

// MockNotifier implements Notifier for testing
type MockNotifier struct {
	name      string
	enabled   bool
	sendError error
	lastMsg   Message
	sendCount int
}

func (m *MockNotifier) Name() string { return m.name }
func (m *MockNotifier) Send(ctx context.Context, msg Message) error {
	m.sendCount++
	m.lastMsg = msg
	return m.sendError
}
func (m *MockNotifier) IsEnabled() bool { return m.enabled }

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	ctx := context.Background()

	n1 := &MockNotifier{name: "mock1", enabled: true}
	n2 := &MockNotifier{name: "mock2", enabled: false}
	n3 := &MockNotifier{name: "mock3", enabled: true, sendError: errors.New("failed")}

	r.Register(n1)
	r.Register(n2)
	r.Register(n3)

	t.Run("List", func(t *testing.T) {
		list := r.List()
		if len(list) != 3 {
			t.Errorf("expected 3 notifiers, got %d", len(list))
		}
	})

	t.Run("Get", func(t *testing.T) {
		n, ok := r.Get("mock1")
		if !ok || n.Name() != "mock1" {
			t.Error("failed to get mock1")
		}
	})

	t.Run("NotifyAll", func(t *testing.T) {
		msg := Message{Title: "Test", Body: "Hello", Severity: SeverityInfo}
		errs := r.NotifyAll(ctx, msg)

		// n1 should succeed
		if n1.sendCount != 1 || n1.lastMsg.Title != "Test" {
			t.Error("n1 should have received the message")
		}

		// n2 should be skipped (disabled)
		if n2.sendCount != 0 {
			t.Error("n2 should have been skipped")
		}

		// n3 should fail
		if n3.sendCount != 1 {
			t.Error("n3 should have tried to send")
		}
		if len(errs) != 1 {
			t.Errorf("expected 1 error from NotifyAll, got %d", len(errs))
		}
	})
}
