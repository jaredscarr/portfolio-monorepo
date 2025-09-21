package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulationGates_ShouldSimulateNetworkDelays(t *testing.T) {
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
			name:           "simulation mode on, network delays on",
			simModeEnabled: true,
			flagValue:      true,
			expected:       true,
		},
		{
			name:           "simulation mode on, network delays off",
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
				mockClient.On("GetFlag", "local", "simulate_network_delays").Return(tt.flagValue, nil)
			}
			
			gates := NewSimulationGates(mockClient, "local")
			result := gates.ShouldSimulateNetworkDelays()
			
			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}
