package gates

// SimulationGatesInterface defines the interface for simulation gates
type SimulationGatesInterface interface {
	IsSimulationModeEnabled() bool
	ShouldDisablePublishing() bool
	ShouldSimulateWebhookFailures() bool
	ShouldSimulateNetworkDelays() bool
	ShouldUsePartialFailureMode() bool
	ShouldUseCircuitBreakerDemo() bool
	GetSimulationStatus() map[string]bool
}
