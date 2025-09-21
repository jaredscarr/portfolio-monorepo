package gates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimulationGates_CheckCircuitBreaker(t *testing.T) {
	tests := []struct {
		name           string
		circuitEnabled bool
		circuitState   CircuitBreakerState
		lastFailure    time.Time
		expected       bool
	}{
		{
			name:           "circuit breaker disabled",
			circuitEnabled: false,
			circuitState:   CircuitOpen,
			expected:       false,
		},
		{
			name:           "circuit closed - allow request",
			circuitEnabled: true,
			circuitState:   CircuitClosed,
			expected:       false,
		},
		{
			name:           "circuit open - recent failure - block request",
			circuitEnabled: true,
			circuitState:   CircuitOpen,
			lastFailure:    time.Now().Add(-2 * time.Second),
			expected:       true,
		},
		{
			name:           "circuit open - old failure - allow test request",
			circuitEnabled: true,
			circuitState:   CircuitOpen,
			lastFailure:    time.Now().Add(-6 * time.Second),
			expected:       false,
		},
		{
			name:           "circuit half-open - allow test request",
			circuitEnabled: true,
			circuitState:   CircuitHalfOpen,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockFeatureFlagClient{}
			mockClient.On("GetFlag", "local", "simulation_mode_enabled").Return(true, nil)
			mockClient.On("GetFlag", "local", "circuit_breaker_demo_mode").Return(tt.circuitEnabled, nil)
			
			gates := NewSimulationGates(mockClient, "local")
			
			// Set up circuit breaker state
			gates.circuitBreaker.state = tt.circuitState
			gates.circuitBreaker.lastFailureTime = tt.lastFailure
			
			result := gates.CheckCircuitBreaker()
			
			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestSimulationGates_RecordCircuitBreakerFailure(t *testing.T) {
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
			initialState:   CircuitClosed,
			initialCount:   0,
			expectedState:  CircuitClosed,
			expectedCount:  0,
		},
		{
			name:           "first failure - stay closed",
			circuitEnabled: true,
			initialState:   CircuitClosed,
			initialCount:   0,
			expectedState:  CircuitClosed,
			expectedCount:  1,
		},
		{
			name:           "third failure - trip to open",
			circuitEnabled: true,
			initialState:   CircuitClosed,
			initialCount:   2,
			expectedState:  CircuitOpen,
			expectedCount:  3,
		},
		{
			name:           "failure in half-open - back to open",
			circuitEnabled: true,
			initialState:   CircuitHalfOpen,
			initialCount:   3,
			expectedState:  CircuitOpen,
			expectedCount:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockFeatureFlagClient{}
			mockClient.On("GetFlag", "local", "simulation_mode_enabled").Return(true, nil)
			mockClient.On("GetFlag", "local", "circuit_breaker_demo_mode").Return(tt.circuitEnabled, nil)
			
			gates := NewSimulationGates(mockClient, "local")
			
			// Set up initial circuit breaker state
			gates.circuitBreaker.state = tt.initialState
			gates.circuitBreaker.failureCount = tt.initialCount
			
			gates.RecordCircuitBreakerFailure()
			
			assert.Equal(t, tt.expectedState, gates.circuitBreaker.state)
			assert.Equal(t, tt.expectedCount, gates.circuitBreaker.failureCount)
			mockClient.AssertExpectations(t)
		})
	}
}
