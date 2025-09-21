package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulationGates_ShouldUsePartialFailureMode(t *testing.T) {
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
			name:           "simulation mode on, partial failure on",
			simModeEnabled: true,
			flagValue:      true,
			expected:       true,
		},
		{
			name:           "simulation mode on, partial failure off",
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
				mockClient.On("GetFlag", "local", "partial_failure_mode").Return(tt.flagValue, nil)
			}
			
			gates := NewSimulationGates(mockClient, "local")
			result := gates.ShouldUsePartialFailureMode()
			
			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}
