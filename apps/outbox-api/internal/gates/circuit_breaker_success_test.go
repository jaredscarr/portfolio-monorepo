package gates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulationGates_RecordCircuitBreakerSuccess(t *testing.T) {
	tests := []struct {
		name           string
		circuitEnabled bool
		initialState   CircuitBreakerState
		initialCount   int
		expectedState  CircuitBreakerState
		expectedCount  int
	}{
		{
			name:           "circuit breaker disabled - no change",
			circuitEnabled: false,
			initialState:   CircuitHalfOpen,
			initialCount:   3,
			expectedState:  CircuitHalfOpen,
			expectedCount:  3,
		},
		{
			name:           "success in closed state - no change",
			circuitEnabled: true,
			initialState:   CircuitClosed,
			initialCount:   1,
			expectedState:  CircuitClosed,
			expectedCount:  1,
		},
		{
			name:           "success in half-open - recover to closed",
			circuitEnabled: true,
			initialState:   CircuitHalfOpen,
			initialCount:   3,
			expectedState:  CircuitClosed,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockFeatureFlagClient{}
			mockClient.On("GetFlag", "local", "simulation_mode_enabled").Return(true, nil)
			mockClient.On("GetFlag", "local", "circuit_breaker_demo_mode").Return(tt.circuitEnabled, nil)
			
			gates := NewSimulationGates(mockClient, "local")
			
			gates.circuitBreaker.state = tt.initialState
			gates.circuitBreaker.failureCount = tt.initialCount
			
			gates.RecordCircuitBreakerSuccess()
			
			assert.Equal(t, tt.expectedState, gates.circuitBreaker.state)
			assert.Equal(t, tt.expectedCount, gates.circuitBreaker.failureCount)
			mockClient.AssertExpectations(t)
		})
	}
}
