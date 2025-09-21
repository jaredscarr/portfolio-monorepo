package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFeatureFlagClient is a mock implementation of FeatureFlagClient
type MockFeatureFlagClient struct {
	mock.Mock
}

func (m *MockFeatureFlagClient) GetFlag(env, key string) (bool, error) {
	args := m.Called(env, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockFeatureFlagClient) GetAllFlags(env string) (map[string]bool, error) {
	args := m.Called(env)
	return args.Get(0).(map[string]bool), args.Error(1)
}

func TestSimulationGates_IsSimulationModeEnabled(t *testing.T) {
	tests := []struct {
		name     string
		flagValue bool
		flagError error
		expected  bool
	}{
		{
			name:      "simulation mode enabled",
			flagValue: true,
			flagError: nil,
			expected:  true,
		},
		{
			name:      "simulation mode disabled", 
			flagValue: false,
			flagError: nil,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockFeatureFlagClient{}
			mockClient.On("GetFlag", "local", "simulation_mode_enabled").Return(tt.flagValue, tt.flagError)
			
			gates := NewSimulationGates(mockClient, "local")
			result := gates.IsSimulationModeEnabled()
			
			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}
 
func TestSimulationGates_ShouldDisablePublishing(t *testing.T) {
	tests := []struct {
		name           string
		simModeEnabled bool
		flagValue      bool
		expected       bool
	}{
		{
			name:           "simulation mode off - should return false",
			simModeEnabled: false,
			flagValue:      true,
			expected:       false,
		},
		{
			name:           "simulation mode on, disable publishing on",
			simModeEnabled: true,
			flagValue:      true,
			expected:       true,
		},
		{
			name:           "simulation mode on, disable publishing off",
			simModeEnabled: true,
			flagValue:      false,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockFeatureFlagClient{}
			mockClient.On("GetFlag", "local", "simulation_mode_enabled").Return(tt.simModeEnabled, nil)
			
			if tt.simModeEnabled {
				mockClient.On("GetFlag", "local", "disable_publishing").Return(tt.flagValue, nil)
			}
			
			gates := NewSimulationGates(mockClient, "local")
			result := gates.ShouldDisablePublishing()
			
			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestSimulationGates_ShouldSimulateWebhookFailures(t *testing.T) {
	tests := []struct {
		name           string
		simModeEnabled bool
		flagValue      bool
		expected       bool
	}{
		{
			name:           "simulation mode off - should return false",
			simModeEnabled: false,
			flagValue:      true,
			expected:       false,
		},
		{
			name:           "simulation mode on, force failures on",
			simModeEnabled: true,
			flagValue:      true,
			expected:       true,
		},
		{
			name:           "simulation mode on, force failures off",
			simModeEnabled: true,
			flagValue:      false,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockFeatureFlagClient{}
			mockClient.On("GetFlag", "local", "simulation_mode_enabled").Return(tt.simModeEnabled, nil)
			
			if tt.simModeEnabled {
				mockClient.On("GetFlag", "local", "force_webhook_failures").Return(tt.flagValue, nil)
			}
			
			gates := NewSimulationGates(mockClient, "local")
			result := gates.ShouldSimulateWebhookFailures()
			
			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}
