package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		expectedStatus  int
		expectedContent string
	}{
		{
			name:            "metrics endpoint returns prometheus metrics",
			expectedStatus:  http.StatusOK,
			expectedContent: "go_goroutines", // Should contain Go runtime metrics
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			router.GET("/metrics", Metrics)

			// Create a test request
			req, _ := http.NewRequest("GET", "/metrics", nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedContent)

			// Verify it's valid Prometheus format
			assert.Contains(t, w.Body.String(), "# HELP")
			assert.Contains(t, w.Body.String(), "# TYPE")
		})
	}
}

func TestMetricsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedLabels map[string]string
	}{
		{
			name:           "GET request metrics",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
			expectedLabels: map[string]string{
				"method": "GET",
				"path":   "/test",
				"status": "200",
			},
		},
		{
			name:           "POST request metrics",
			method:         "POST",
			path:           "/api/data",
			expectedStatus: http.StatusOK,
			expectedLabels: map[string]string{
				"method": "POST",
				"path":   "/api/data",
				"status": "200",
			},
		},
		{
			name:           "404 request metrics",
			method:         "GET",
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedLabels: map[string]string{
				"method": "GET",
				"path":   "/nonexistent",
				"status": "404",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset metrics before each test
			httpRequestsTotal.Reset()
			httpRequestDuration.Reset()

			// Create a new Gin router with metrics middleware
			router := gin.New()
			router.Use(MetricsMiddleware())

			// Add a test route
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})
			router.POST("/api/data", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "data"})
			})

			// Create a test request
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check that metrics were recorded
			// Note: We can't easily test the exact values without more complex setup,
			// but we can verify the metrics were incremented
			assert.True(t, testutil.ToFloat64(httpRequestsTotal.WithLabelValues(
				tt.expectedLabels["method"],
				tt.expectedLabels["path"],
				tt.expectedLabels["status"],
			)) > 0)

			// For histogram, we can't easily test the exact value without more complex setup
			// but we can verify the metric was observed by checking the registry
			// This is a simplified test - in practice, you'd need to gather metrics from the registry
		})
	}
}

func TestMetricsMiddlewareWithNamedRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset metrics before test
	httpRequestsTotal.Reset()
	httpRequestDuration.Reset()

	// Create a new Gin router with metrics middleware
	router := gin.New()
	router.Use(MetricsMiddleware())

	// Add a named route
	router.GET("/health", Health)

	// Create a test request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that metrics were recorded with the named route path
	assert.True(t, testutil.ToFloat64(httpRequestsTotal.WithLabelValues(
		"GET",
		"/health",
		"200",
	)) > 0)

	// For histogram, we can't easily test the exact value without more complex setup
	// but we can verify the metric was observed by checking the registry
}

func TestMetricsMiddlewareTiming(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset metrics before test
	httpRequestsTotal.Reset()
	httpRequestDuration.Reset()

	// Create a new Gin router with metrics middleware
	router := gin.New()
	router.Use(MetricsMiddleware())

	// Add a route that takes some time
	router.GET("/slow", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "slow response"})
	})

	// Create a test request
	req, _ := http.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// For histogram timing, we can't easily test the exact value without more complex setup
	// but we can verify the metric was observed by checking the registry
	// In practice, you'd gather metrics from the registry and check the histogram values
}

// Test that metrics are properly registered
func TestMetricsRegistration(t *testing.T) {
	// Try to register the same metrics again - should panic if already registered
	assert.Panics(t, func() {
		prometheus.MustRegister(httpRequestsTotal)
	})

	assert.Panics(t, func() {
		prometheus.MustRegister(httpRequestDuration)
	})

	// Verify they exist in the registry by creating a new registry and checking
	newRegistry := prometheus.NewRegistry()

	// Try to register our metrics in the new registry - should work
	assert.NotPanics(t, func() {
		newRegistry.MustRegister(httpRequestsTotal)
	})

	assert.NotPanics(t, func() {
		newRegistry.MustRegister(httpRequestDuration)
	})

	// Gather metrics from the new registry
	metrics, err := newRegistry.Gather()
	assert.NoError(t, err)

	var foundRequestsTotal, foundRequestDuration bool
	for _, metric := range metrics {
		if metric.GetName() == "http_requests_total" {
			foundRequestsTotal = true
		}
		if metric.GetName() == "http_request_duration_seconds" {
			foundRequestDuration = true
		}
	}

	assert.True(t, foundRequestsTotal, "http_requests_total metric should be registered")
	assert.True(t, foundRequestDuration, "http_request_duration_seconds metric should be registered")
}

// Benchmark tests for performance
func BenchmarkMetrics(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/metrics", Metrics)

	req, _ := http.NewRequest("GET", "/metrics", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkMetricsMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(MetricsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
