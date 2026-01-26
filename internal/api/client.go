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

// extractProjectId handles the complex type of cloudaicompanionProject
func extractProjectId(value interface{}) string {
	if value == nil {
		return ""
	}

	// Case 1: String
	if str, ok := value.(string); ok && str != "" {
		return str
	}

	// Case 2: Object with ID
	if m, ok := value.(map[string]interface{}); ok {
		if id, ok := m["id"].(string); ok && id != "" {
			return id
		}
	}

	return ""
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

	// Try to resolve project ID from various fields
	projectID := response.ProjectID
	if projectID == "" {
		projectID = extractProjectId(response.CloudAICompanionProject)
	}

	// Store project ID if available
	if projectID != "" {
		c.SetProjectID(projectID)
		response.ProjectID = projectID // Ensure it's set in the response struct too
	}

	return &response, nil
}

// OnboardUser attempts to onboard the user to get a project ID
func (c *Client) OnboardUser(tierID string) (string, error) {
	request := models.OnboardUserRequest{
		TierID:   tierID,
		Metadata: models.GetDefaultMetadata(),
	}

	for attempt := 1; attempt <= 5; attempt++ {
		responseData, err := c.doRequest("POST", "/v1internal:onboardUser", request)
		if err != nil {
			return "", fmt.Errorf("onboardUser failed: %w", err)
		}

		var response models.OnboardUserResponse
		if err := json.Unmarshal(responseData, &response); err != nil {
			return "", fmt.Errorf("failed to parse onboard response: %w", err)
		}

		if response.Done {
			projectID := extractProjectId(response.Response.CloudAICompanionProject)
			if projectID != "" {
				return projectID, nil
			}
			// Done but no project ID?
			return "", nil
		}

		// Wait before retry
		time.Sleep(2 * time.Second)
	}

	return "", fmt.Errorf("onboarding timed out")
}

// ResolveProjectID implements the full logic to get a project ID
func (c *Client) ResolveProjectID() (string, error) {
	// Step 1: Call loadCodeAssist
	resp, err := c.LoadCodeAssist()
	if err != nil {
		return "", err
	}

	// Step 2: Check if we already got it
	if resp.ProjectID != "" {
		return resp.ProjectID, nil
	}

	// Step 3: Determine Tier
	var tierID string
	if resp.PaidTier != nil && resp.PaidTier.ID != "" {
		tierID = resp.PaidTier.ID
	} else if resp.CurrentTier != nil && resp.CurrentTier.ID != "" {
		tierID = resp.CurrentTier.ID
	} else {
		// Pick from allowed tiers
		if len(resp.AllowedTiers) > 0 {
			// Find default
			for _, t := range resp.AllowedTiers {
				if t.IsDefault && t.ID != "" {
					tierID = t.ID
					break
				}
			}
			// Or first
			if tierID == "" && resp.AllowedTiers[0].ID != "" {
				tierID = resp.AllowedTiers[0].ID
			}

			// Fallback
			if tierID == "" {
				tierID = "LEGACY"
			}
		}
	}

	if tierID == "" {
		return "", fmt.Errorf("cannot determine tier for onboarding")
	}

	// Step 4: Onboard
	projectID, err := c.OnboardUser(tierID)
	if err != nil {
		return "", err
	}

	if projectID != "" {
		c.SetProjectID(projectID)
	}

	return projectID, nil
}

// FetchAvailableModels retrieves available models with quota information
func (c *Client) FetchAvailableModels() (*models.FetchAvailableModelsResponse, error) {
	if err := c.EnsureAuthenticated(); err != nil {
		return nil, err
	}

	// Ensure we have a project ID first (via ResolveProjectID)
	// We don't call it here to avoid recursion loop, but assumed it's called before or handled

	request := models.FetchAvailableModelsRequest{}

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
	// First, resolve project ID (this handles onboarding if needed)
	_, err := c.ResolveProjectID()
	if err != nil {
		fmt.Printf("DEBUG: ResolveProjectID failed: %v\n", err)
		// Not a fatal error, continue without project ID
	} else {
		// Log success for debugging
		// fmt.Printf("DEBUG: Project ID resolved: %s\n", projectID)
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
