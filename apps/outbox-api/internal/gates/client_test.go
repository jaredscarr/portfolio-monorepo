package gates

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPFeatureFlagClient_GetFlag(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedValue  bool
		expectError    bool
	}{
		{
			name:           "successful flag retrieval - true",
			responseStatus: http.StatusOK,
			responseBody:   `{"key": "test_flag", "enabled": true}`,
			expectedValue:  true,
			expectError:    false,
		},
		{
			name:           "successful flag retrieval - false",
			responseStatus: http.StatusOK,
			responseBody:   `{"key": "test_flag", "enabled": false}`,
			expectedValue:  false,
			expectError:    false,
		},
		{
			name:           "flag not found",
			responseStatus: http.StatusNotFound,
			responseBody:   `{"error": "flag not found"}`,
			expectedValue:  false,
			expectError:    true,
		},
		{
			name:           "server error",
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "internal error"}`,
			expectedValue:  false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/flags/test_flag", r.URL.Path)
				assert.Equal(t, "local", r.URL.Query().Get("env"))
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewHTTPFeatureFlagClient(server.URL)
			value, err := client.GetFlag("local", "test_flag")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}

func TestHTTPFeatureFlagClient_GetAllFlags(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedFlags  map[string]bool
		expectError    bool
	}{
		{
			name:           "successful flags retrieval",
			responseStatus: http.StatusOK,
			responseBody:   `{"simulation_mode_enabled": true, "disable_publishing": false}`,
			expectedFlags:  map[string]bool{"simulation_mode_enabled": true, "disable_publishing": false},
			expectError:    false,
		},
		{
			name:           "empty flags",
			responseStatus: http.StatusOK,
			responseBody:   `{}`,
			expectedFlags:  map[string]bool{},
			expectError:    false,
		},
		{
			name:           "server error",
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "internal error"}`,
			expectedFlags:  nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/flags", r.URL.Path)
				assert.Equal(t, "local", r.URL.Query().Get("env"))
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewHTTPFeatureFlagClient(server.URL)
			flags, err := client.GetAllFlags("local")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFlags, flags)
			}
		})
	}
}
