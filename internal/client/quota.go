package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	// Default paths relative to user home
	OAuthCredsPath = ".gemini/oauth_creds.json"
)

type OAuthCreds struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"` // Assuming unix timestamp or similar
}

type Client struct {
	baseURL    string
	token      string
	csrfToken  string
	httpClient *http.Client
}

type ModelQuota struct {
	Name       string    `json:"name"`
	QuotaUsed  int64     `json:"quotaUsed"`
	QuotaLimit int64     `json:"quotaLimit"`
	ResetAt    time.Time `json:"resetAt"`
}

type UserStatusResponse struct {
	Models []ModelQuota `json:"models"`
}

func NewClient(port int, csrfToken string) (*Client, error) {
	token, err := readAuthToken()
	if err != nil {
		return nil, fmt.Errorf("failed to read auth token: %w", err)
	}

	return &Client{
		baseURL:    fmt.Sprintf("http://localhost:%d", port),
		token:      token,
		csrfToken:  csrfToken,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func readAuthToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	credsPath := filepath.Join(home, OAuthCredsPath)
	data, err := os.ReadFile(credsPath)
	if err != nil {
		return "", err
	}

	var creds OAuthCreds
	if err := json.Unmarshal(data, &creds); err != nil {
		return "", err
	}

	if creds.AccessToken == "" {
		return "", fmt.Errorf("access token is empty")
	}

	return creds.AccessToken, nil
}

func (c *Client) GetUserStatus() (*UserStatusResponse, error) {
	url := fmt.Sprintf("%s/GetUserStatus", c.baseURL)
	if c.csrfToken != "" {
		url = fmt.Sprintf("%s?csrf_token=%s", url, c.csrfToken)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	if c.csrfToken != "" {
		req.Header.Set("X-CSRF-Token", c.csrfToken)
		req.Header.Set("X-XSRF-Token", c.csrfToken)
		req.Header.Set("X-Csrf-Token", c.csrfToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var statusResp UserStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResp, nil
}
