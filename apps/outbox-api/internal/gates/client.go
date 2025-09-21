package gates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FeatureFlagClient defines the interface for fetching feature flags
type FeatureFlagClient interface {
	GetFlag(env, key string) (bool, error)
	GetAllFlags(env string) (map[string]bool, error)
}

// HTTPFeatureFlagClient implements FeatureFlagClient using HTTP calls to the feature-flags-api
type HTTPFeatureFlagClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPFeatureFlagClient creates a new HTTP-based feature flag client
func NewHTTPFeatureFlagClient(baseURL string) *HTTPFeatureFlagClient {
	return &HTTPFeatureFlagClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetFlag retrieves a single feature flag value
func (c *HTTPFeatureFlagClient) GetFlag(env, key string) (bool, error) {
	url := fmt.Sprintf("%s/flags/%s?env=%s", c.baseURL, key, env)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to fetch flag %s: %w", key, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, fmt.Errorf("flag %s not found", key)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("feature flags API returned status %d", resp.StatusCode)
	}

	var flagResponse struct {
		Key     string `json:"key"`
		Enabled bool   `json:"enabled"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&flagResponse); err != nil {
		return false, fmt.Errorf("failed to decode flag response: %w", err)
	}

	return flagResponse.Enabled, nil
}

// GetAllFlags retrieves all feature flags for an environment
func (c *HTTPFeatureFlagClient) GetAllFlags(env string) (map[string]bool, error) {
	url := fmt.Sprintf("%s/flags?env=%s", c.baseURL, env)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feature flags API returned status %d", resp.StatusCode)
	}

	var flags map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&flags); err != nil {
		return nil, fmt.Errorf("failed to decode flags response: %w", err)
	}

	return flags, nil
}
