package gates

import (
	"log"
	"sync"
	"time"
)

type CircuitBreakerState int

const (
	CircuitClosed CircuitBreakerState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreakerData struct {
	state           CircuitBreakerState
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.RWMutex
}

type SimulationGates struct {
	flagsClient   FeatureFlagClient
	environment   string
	circuitBreaker *CircuitBreakerData
}

func NewSimulationGates(flagsClient FeatureFlagClient, environment string) *SimulationGates {
	return &SimulationGates{
		flagsClient: flagsClient,
		environment: environment,
		circuitBreaker: &CircuitBreakerData{
			state: CircuitClosed,
		},
	}
}

// IsSimulationModeEnabled checks if simulation mode is enabled
func (g *SimulationGates) IsSimulationModeEnabled() bool {
	enabled, err := g.flagsClient.GetFlag(g.environment, "simulation_mode_enabled")
	if err != nil {
		log.Printf("Warning: failed to get simulation_mode_enabled flag: %v", err)
		return false
	}
	return enabled
}

// ShouldSimulateWebhookFailures determines if webhook calls should be forced to fail
func (g *SimulationGates) ShouldSimulateWebhookFailures() bool {
	if !g.IsSimulationModeEnabled() {
		return false
	}
	
	enabled, err := g.flagsClient.GetFlag(g.environment, "force_webhook_failures")
	if err != nil {
		log.Printf("Warning: failed to get force_webhook_failures flag: %v", err)
		return false
	}
	return enabled
}

// ShouldDisablePublishing determines if publishing should be completely disabled
func (g *SimulationGates) ShouldDisablePublishing() bool {
	if !g.IsSimulationModeEnabled() {
		return false
	}
	
	enabled, err := g.flagsClient.GetFlag(g.environment, "disable_publishing")
	if err != nil {
		log.Printf("Warning: failed to get disable_publishing flag: %v", err)
		return false
	}
	return enabled
}

// ShouldUseCircuitBreakerDemo determines if circuit breaker demo mode is active
func (g *SimulationGates) ShouldUseCircuitBreakerDemo() bool {
	if !g.IsSimulationModeEnabled() {
		return false
	}
	
	enabled, err := g.flagsClient.GetFlag(g.environment, "circuit_breaker_demo_mode")
	if err != nil {
		log.Printf("Warning: failed to get circuit_breaker_demo_mode flag: %v", err)
		return false
	}
	return enabled
}

// ShouldUsePartialFailureMode determines if some events should succeed and others fail
func (g *SimulationGates) ShouldUsePartialFailureMode() bool {
	if !g.IsSimulationModeEnabled() {
		return false
	}
	
	enabled, err := g.flagsClient.GetFlag(g.environment, "partial_failure_mode")
	if err != nil {
		log.Printf("Warning: failed to get partial_failure_mode flag: %v", err)
		return false
	}
	return enabled
}

// ShouldSimulateNetworkDelays determines if artificial delays should be added
func (g *SimulationGates) ShouldSimulateNetworkDelays() bool {
	if !g.IsSimulationModeEnabled() {
		return false
	}
	
	enabled, err := g.flagsClient.GetFlag(g.environment, "simulate_network_delays")
	if err != nil {
		log.Printf("Warning: failed to get simulate_network_delays flag: %v", err)
		return false
	}
	return enabled
}

// CheckCircuitBreaker checks if circuit breaker should block the request
func (g *SimulationGates) CheckCircuitBreaker() bool {
	if !g.ShouldUseCircuitBreakerDemo() {
		return false // Circuit breaker not active
	}
	
	g.circuitBreaker.mutex.RLock()
	state := g.circuitBreaker.state
	lastFailureTime := g.circuitBreaker.lastFailureTime
	g.circuitBreaker.mutex.RUnlock()
	
	switch state {
	case CircuitClosed:
		return false // Allow request
	case CircuitOpen:
		// Check if timeout period has passed (5 seconds)
		if time.Since(lastFailureTime) > 5*time.Second {
			// Move to half-open state
			g.circuitBreaker.mutex.Lock()
			g.circuitBreaker.state = CircuitHalfOpen
			g.circuitBreaker.mutex.Unlock()
			log.Printf("Circuit breaker: OPEN → HALF-OPEN (testing)")
			return false // Allow one test request
		}
		return true // Block request
	case CircuitHalfOpen:
		return false // Allow test request
	default:
		return false
	}
}

// RecordCircuitBreakerSuccess records a successful request
func (g *SimulationGates) RecordCircuitBreakerSuccess() {
	if !g.ShouldUseCircuitBreakerDemo() {
		return
	}
	
	g.circuitBreaker.mutex.Lock()
	defer g.circuitBreaker.mutex.Unlock()
	
	if g.circuitBreaker.state == CircuitHalfOpen {
		g.circuitBreaker.state = CircuitClosed
		g.circuitBreaker.failureCount = 0
		log.Printf("Circuit breaker: HALF-OPEN → CLOSED (recovered)")
	}
}

// RecordCircuitBreakerFailure records a failed request
func (g *SimulationGates) RecordCircuitBreakerFailure() {
	if !g.ShouldUseCircuitBreakerDemo() {
		return
	}
	
	g.circuitBreaker.mutex.Lock()
	defer g.circuitBreaker.mutex.Unlock()
	
	g.circuitBreaker.failureCount++
	g.circuitBreaker.lastFailureTime = time.Now()
	
	// Trip circuit after 3 failures
	if g.circuitBreaker.failureCount >= 3 && g.circuitBreaker.state == CircuitClosed {
		g.circuitBreaker.state = CircuitOpen
		log.Printf("Circuit breaker: CLOSED → OPEN (tripped after %d failures)", g.circuitBreaker.failureCount)
	} else if g.circuitBreaker.state == CircuitHalfOpen {
		g.circuitBreaker.state = CircuitOpen
		log.Printf("Circuit breaker: HALF-OPEN → OPEN (test failed)")
	}
}

// GetSimulationStatus returns a summary of current simulation settings
func (g *SimulationGates) GetSimulationStatus() map[string]interface{} {
	g.circuitBreaker.mutex.RLock()
	circuitState := g.circuitBreaker.state
	failureCount := g.circuitBreaker.failureCount
	lastFailureTime := g.circuitBreaker.lastFailureTime
	g.circuitBreaker.mutex.RUnlock()
	
	var circuitStateStr string
	switch circuitState {
	case CircuitClosed:
		circuitStateStr = "CLOSED"
	case CircuitOpen:
		circuitStateStr = "OPEN"
	case CircuitHalfOpen:
		circuitStateStr = "HALF-OPEN"
	default:
		circuitStateStr = "UNKNOWN"
	}
	
	return map[string]interface{}{
		"simulation_mode_enabled":    g.IsSimulationModeEnabled(),
		"force_webhook_failures":     g.ShouldSimulateWebhookFailures(),
		"disable_publishing":         g.ShouldDisablePublishing(),
		"circuit_breaker_demo_mode":  g.ShouldUseCircuitBreakerDemo(),
		"partial_failure_mode":       g.ShouldUsePartialFailureMode(),
		"simulate_network_delays":    g.ShouldSimulateNetworkDelays(),
		"circuit_breaker_state":      circuitStateStr,
		"circuit_failure_count":      failureCount,
		"circuit_last_failure":       lastFailureTime,
	}
}
