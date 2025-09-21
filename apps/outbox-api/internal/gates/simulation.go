package gates

import (
	"log"
)

// SimulationGates encapsulates feature flag logic for simulation behaviors
type SimulationGates struct {
	flagsClient FeatureFlagClient
	environment string
}

// NewSimulationGates creates a new SimulationGates instance
func NewSimulationGates(flagsClient FeatureFlagClient, environment string) *SimulationGates {
	return &SimulationGates{
		flagsClient: flagsClient,
		environment: environment,
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

// GetSimulationStatus returns a summary of current simulation settings
func (g *SimulationGates) GetSimulationStatus() map[string]bool {
	return map[string]bool{
		"simulation_mode_enabled":    g.IsSimulationModeEnabled(),
		"force_webhook_failures":     g.ShouldSimulateWebhookFailures(),
		"disable_publishing":         g.ShouldDisablePublishing(),
		"circuit_breaker_demo_mode":  g.ShouldUseCircuitBreakerDemo(),
		"partial_failure_mode":       g.ShouldUsePartialFailureMode(),
		"simulate_network_delays":    g.ShouldSimulateNetworkDelays(),
	}
}
