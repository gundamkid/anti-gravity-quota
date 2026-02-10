package notify

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// RoundTripFunc is a mock transport
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestTelegramNotifier_Send(t *testing.T) {
	ctx := context.Background()
	msg := Message{
		Title:    "Test Title",
		Body:     "Test Body",
		Severity: SeverityCritical,
	}

	t.Run("Success", func(t *testing.T) {
		client := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				// Check request
				if !strings.Contains(req.URL.Path, "/botfake-token/sendMessage") {
					t.Errorf("unexpected URL: %s", req.URL.Path)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
					Header:     make(http.Header),
				}
			}),
		}

		tn := &TelegramNotifier{
			token:  "fake-token",
			chatID: "fake-chat",
			client: client,
		}

		err := tn.Send(ctx, msg)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("API Error", func(t *testing.T) {
		client := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ok": false, "description": "bad request"}`)),
					Header:     make(http.Header),
				}
			}),
		}

		tn := &TelegramNotifier{
			token:  "fake-token",
			chatID: "fake-chat",
			client: client,
		}

		err := tn.Send(ctx, msg)
		if err == nil || !strings.Contains(err.Error(), "bad request") {
			t.Errorf("expected bad request error, got %v", err)
		}
	})

	t.Run("Rate Limiting", func(t *testing.T) {
		client := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
					Header:     make(http.Header),
				}
			}),
		}

		tn := &TelegramNotifier{
			token:  "fake-token",
			chatID: "fake-chat",
			client: client,
		}

		// Fill the rate limit
		for i := 0; i < 10; i++ {
			if err := tn.Send(ctx, msg); err != nil {
				t.Fatalf("failed at message %d: %v", i, err)
			}
		}

		// 11th message should fail
		err := tn.Send(ctx, msg)
		if err == nil || !strings.Contains(err.Error(), "rate limit exceeded") {
			t.Errorf("expected rate limit error, got %v", err)
		}
	})
}

func TestTelegramNotifier_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid Token", func(t *testing.T) {
		client := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ok": true}`)),
					Header:     make(http.Header),
				}
			}),
		}

		tn := &TelegramNotifier{
			token:  "valid-token",
			client: client,
		}

		err := tn.Validate(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		client := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{"ok": false, "description": "unauthorized"}`)),
					Header:     make(http.Header),
				}
			}),
		}

		tn := &TelegramNotifier{
			token:  "invalid-token",
			client: client,
		}

		err := tn.Validate(ctx)
		if err == nil || !strings.Contains(err.Error(), "unauthorized") {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})
}
