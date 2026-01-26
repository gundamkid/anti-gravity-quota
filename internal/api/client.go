package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/gundamkid/anti-gravity-quota/internal/models"
)

const (
	// BaseURL is the primary Cloud Code API endpoint
	BaseURL = "https://cloudcode-pa.googleapis.com"

	// BackupURL is the backup Cloud Code API endpoint
	BackupURL = "https://daily-cloudcode-pa.sandbox.googleapis.com"

	// UserAgent for API requests
	UserAgent = "antigravity"

	// MaxRetries for failed requests
	MaxRetries = 3

	// RetryDelay initial delay between retries
	RetryDelay = 1 * time.Second
)

// Client represents a Cloud Code API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
	projectID  string
}

// NewClient creates a new Cloud Code API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: BaseURL,
	}
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetProjectID sets the project ID
func (c *Client) SetProjectID(projectID string) {
	c.projectID = projectID
}

// GetProjectID returns the current project ID
func (c *Client) GetProjectID() string {
	return c.projectID
}

// doRequest performs an HTTP request with authentication headers and retry logic
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := RetryDelay * time.Duration(1<<uint(attempt-1))
			time.Sleep(delay)
		}

		// Marshal request body
		var bodyReader io.Reader
		if body != nil {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request: %w", err)
			}
			bodyReader = bytes.NewReader(jsonData)
		}

		// Create request
		url := c.baseURL + endpoint
		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", UserAgent)
		if c.token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		}

		// Perform request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		// Read response
		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		// Check status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return responseBody, nil
		}

		// Handle errors
		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("unauthorized: token may be invalid or expired")
		}

		if resp.StatusCode == 429 {
			// Rate limited, retry with backoff
			lastErr = fmt.Errorf("rate limited (attempt %d/%d)", attempt+1, MaxRetries+1)
			continue
		}

		if resp.StatusCode >= 500 {
			// Server error, retry
			lastErr = fmt.Errorf("server error %d (attempt %d/%d)", resp.StatusCode, attempt+1, MaxRetries+1)
			continue
		}

		// Client error, don't retry
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", MaxRetries+1, lastErr)
}

// EnsureAuthenticated ensures the client has a valid token
func (c *Client) EnsureAuthenticated() error {
	if c.token != "" {
		return nil
	}

	// Get valid token (will auto-refresh if needed)
	token, err := auth.GetValidToken(auth.GetOAuthConfig())
	if err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	c.SetToken(token)
	return nil
}

// LoadCodeAssist loads code assist status and retrieves project ID
func (c *Client) LoadCodeAssist() (*models.LoadCodeAssistResponse, error) {
	if err := c.EnsureAuthenticated(); err != nil {
		return nil, err
	}

	request := models.LoadCodeAssistRequest{
		Metadata: models.GetDefaultMetadata(),
	}

	responseData, err := c.doRequest("POST", "/v1internal:loadCodeAssist", request)
	if err != nil {
		return nil, fmt.Errorf("loadCodeAssist failed: %w", err)
	}

	var response models.LoadCodeAssistResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Store project ID if available
	if response.ProjectID != "" {
		c.SetProjectID(response.ProjectID)
	}

	return &response, nil
}

// FetchAvailableModels retrieves available models with quota information
func (c *Client) FetchAvailableModels() (*models.FetchAvailableModelsResponse, error) {
	if err := c.EnsureAuthenticated(); err != nil {
		return nil, err
	}

	request := models.FetchAvailableModelsRequest{
		Metadata: models.GetDefaultMetadata(),
	}

	// Include project ID if available
	if c.projectID != "" {
		request.Project = c.projectID
	}

	responseData, err := c.doRequest("POST", "/v1internal:fetchAvailableModels", request)
	if err != nil {
		return nil, fmt.Errorf("fetchAvailableModels failed: %w", err)
	}

	var response models.FetchAvailableModelsResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetQuotaInfo retrieves complete quota information for all models
func (c *Client) GetQuotaInfo() (*models.QuotaSummary, error) {
	// First, load code assist to get project ID
	_, err := c.LoadCodeAssist()
	if err != nil {
		// Not a fatal error, continue without project ID
		// Some accounts may not have a project ID
	}

	// Fetch available models
	modelsResp, err := c.FetchAvailableModels()
	if err != nil {
		return nil, err
	}

	// Convert to QuotaSummary
	quotaSummary := &models.QuotaSummary{
		ProjectID:      c.projectID,
		DefaultModelID: modelsResp.DefaultAgentModel,
		FetchedAt:      time.Now(),
		Models:         make([]models.ModelQuota, 0, len(modelsResp.Models)),
	}

	// Get email from token
	token, err := auth.LoadToken()
	if err == nil {
		quotaSummary.Email = token.Email
	}

	// Convert models to ModelQuota
	for modelID, model := range modelsResp.Models {
		quotaSummary.Models = append(quotaSummary.Models, model.ToModelQuota(modelID))
	}

	return quotaSummary, nil
}
