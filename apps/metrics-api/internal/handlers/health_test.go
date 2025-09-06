package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:           "health check returns ok",
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"status": "ok"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			router.GET("/health", Health)

			// Create a test request
			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare the response body
			var actualResponse gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}

func TestReady(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:           "ready check returns ready",
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"status": "ready"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			router.GET("/ready", Ready)

			// Create a test request
			req, _ := http.NewRequest("GET", "/ready", nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare the response body
			var actualResponse gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}

// Test both health endpoints together
func TestHealthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", Health)
	router.GET("/ready", Ready)

	tests := []struct {
		name           string
		endpoint       string
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:           "health endpoint",
			endpoint:       "/health",
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"status": "ok"},
		},
		{
			name:           "ready endpoint",
			endpoint:       "/ready",
			expectedStatus: http.StatusOK,
			expectedBody:   gin.H{"status": "ready"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.endpoint, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actualResponse gin.H
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}

// Benchmark tests for performance
func BenchmarkHealth(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", Health)

	req, _ := http.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkReady(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/ready", Ready)

	req, _ := http.NewRequest("GET", "/ready", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
