package models

import "time"

// ModelQuota represents quota information for a single model
type ModelQuota struct {
	ModelID           string
	DisplayName       string
	Label             string
	Provider          string
	RemainingFraction float64
	ResetTime         time.Time
	IsExhausted       bool
}

// QuotaSummary represents the complete quota information
type QuotaSummary struct {
	Email            string
	ProjectID        string
	Models           []ModelQuota
	DefaultModelID   string
	FetchedAt        time.Time
}

// LoadCodeAssistRequest represents the request to load code assist
type LoadCodeAssistRequest struct {
	Metadata Metadata `json:"metadata"`
}

// LoadCodeAssistResponse represents the response from loadCodeAssist endpoint
type LoadCodeAssistResponse struct {
	ProjectID string `json:"projectId,omitempty"`
	Status    string `json:"status,omitempty"`
}

// FetchAvailableModelsRequest represents the request to fetch models
type FetchAvailableModelsRequest struct {
	Project  string   `json:"project,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// FetchAvailableModelsResponse represents the response from fetchAvailableModels endpoint
type FetchAvailableModelsResponse struct {
	Models            map[string]Model `json:"models"`
	DefaultAgentModel string           `json:"defaultAgentModelId,omitempty"`
}

// Model represents a single model with quota information
type Model struct {
	DisplayName   string         `json:"displayName"`
	Model         string         `json:"model"`
	Label         string         `json:"label"`
	QuotaInfo     ModelQuotaInfo `json:"quotaInfo"`
	ModelProvider string         `json:"modelProvider"`
}

// ModelQuotaInfo represents quota details for a model
type ModelQuotaInfo struct {
	RemainingFraction float64   `json:"remainingFraction"`
	ResetTime         time.Time `json:"resetTime"`
	IsExhausted       bool      `json:"isExhausted"`
}

// Metadata represents request metadata
type Metadata struct {
	IDEType    string `json:"ideType"`
	Platform   string `json:"platform"`
	PluginType string `json:"pluginType"`
}

// GetDefaultMetadata returns the default metadata for API requests
func GetDefaultMetadata() Metadata {
	return Metadata{
		IDEType:    "ANTIGRAVITY",
		Platform:   "PLATFORM_UNSPECIFIED",
		PluginType: "GEMINI",
	}
}

// ToModelQuota converts a Model to ModelQuota
func (m Model) ToModelQuota(modelID string) ModelQuota {
	return ModelQuota{
		ModelID:           modelID,
		DisplayName:       m.DisplayName,
		Label:             m.Label,
		Provider:          m.ModelProvider,
		RemainingFraction: m.QuotaInfo.RemainingFraction,
		ResetTime:         m.QuotaInfo.ResetTime,
		IsExhausted:       m.QuotaInfo.IsExhausted,
	}
}

// GetRemainingPercentage returns the remaining quota as a percentage (0-100)
func (q ModelQuota) GetRemainingPercentage() int {
	return int(q.RemainingFraction * 100)
}

// GetTimeUntilReset returns the duration until quota reset
func (q ModelQuota) GetTimeUntilReset() time.Duration {
	return time.Until(q.ResetTime)
}

// GetStatusString returns a human-readable status string
func (q ModelQuota) GetStatusString() string {
	if q.IsExhausted {
		return "EMPTY"
	}
	if q.RemainingFraction <= 0.1 {
		return "LOW"
	}
	return "OK"
}
