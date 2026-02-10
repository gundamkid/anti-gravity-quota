package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TelegramNotifier implements the Notifier interface for Telegram
type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
	mu     sync.Mutex
	// entries tracks timestamps of messages sent in the last minute for rate limiting
	entries []time.Time
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *TelegramNotifier) Name() string {
	return "telegram"
}

func (t *TelegramNotifier) IsEnabled() bool {
	return t.token != "" && t.chatID != ""
}

// TelegramResponse represents a basic response from Telegram API
type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

// Send sends a message to the configured Telegram chat with rate limiting
func (t *TelegramNotifier) Send(ctx context.Context, msg Message) error {
	if !t.IsEnabled() {
		return fmt.Errorf("telegram notifier not configured")
	}

	// Rate limiting: max 10 messages/minute
	t.mu.Lock()
	now := time.Now()
	// Prune old entries
	var recent []time.Time
	oneMinuteAgo := now.Add(-1 * time.Minute)
	for _, entry := range t.entries {
		if entry.After(oneMinuteAgo) {
			recent = append(recent, entry)
		}
	}
	t.entries = recent

	if len(t.entries) >= 10 {
		t.mu.Unlock()
		return fmt.Errorf("telegram rate limit exceeded (max 10 msgs/min)")
	}
	t.entries = append(t.entries, now)
	t.mu.Unlock()

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	// Format message with MarkdownV2-friendly style
	// severity emoji
	emoji := "ℹ️"
	switch msg.Severity {
	case SeverityWarning:
		emoji = "⚠️"
	case SeverityCritical:
		emoji = "🚨"
	case SeverityRecovery:
		emoji = "✅"
	}

	text := fmt.Sprintf("%s *%s*", emoji, msg.Title)
	text += "\n\n" + msg.Body

	body, err := json.Marshal(map[string]string{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "Markdown", // Using Markdown (v1) for simplicity as per requirements
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var tgResp TelegramResponse
		if err := json.NewDecoder(resp.Body).Decode(&tgResp); err == nil {
			return fmt.Errorf("telegram API error (%d): %s", resp.StatusCode, tgResp.Description)
		}
		return fmt.Errorf("telegram API error: status code %d", resp.StatusCode)
	}

	return nil
}

// Validate tests the telegram bot token using getMe API
func (t *TelegramNotifier) Validate(ctx context.Context) error {
	if t.token == "" {
		return fmt.Errorf("bot token is empty")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", t.token)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var tgResp TelegramResponse
		if err := json.NewDecoder(resp.Body).Decode(&tgResp); err == nil {
			return fmt.Errorf("invalid token: %s", tgResp.Description)
		}
		return fmt.Errorf("invalid token: status code %d", resp.StatusCode)
	}

	return nil
}
