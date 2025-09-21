package gates

// SimulationGatesInterface defines the interface for simulation gates
type SimulationGatesInterface interface {
	IsSimulationModeEnabled() bool
	ShouldDisablePublishing() bool
	ShouldSimulateWebhookFailures() bool
	ShouldSimulateNetworkDelays() bool
	ShouldUsePartialFailureMode() bool
	ShouldUseCircuitBreakerDemo() bool
	CheckCircuitBreaker() bool
	RecordCircuitBreakerSuccess()
	RecordCircuitBreakerFailure()
	GetSimulationStatus() map[string]interface{}
}
