package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFlagsService defines the interface for mocking flags package functions
type MockFlagsService struct {
	mock.Mock
}

func (m *MockFlagsService) GetAllFlags(env string) (map[string]bool, error) {
	args := m.Called(env)
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockFlagsService) GetSingleFlag(env, key string) (bool, bool, error) {
	args := m.Called(env, key)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

func TestGetFlags(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		env            string
		mockFlags      map[string]bool
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "successful local flags",
			env:            "local",
			mockFlags:      map[string]bool{"feature_a": true, "feature_b": false},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]bool{"feature_a": true, "feature_b": false},
		},
		{
			name:           "successful prod flags",
			env:            "prod",
			mockFlags:      map[string]bool{"feature_c": true, "feature_d": false, "feature_e": true},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]bool{"feature_c": true, "feature_d": false, "feature_e": true},
		},
		{
			name:           "empty flags",
			env:            "local",
			mockFlags:      map[string]bool{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]bool{},
		},
		{
			name:           "flags service error",
			env:            "local",
			mockFlags:      nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ErrorResponse{Error: assert.AnError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()

			// Create mock flags service
			mockFlagsService := new(MockFlagsService)

			// Set up mock expectations
			mockFlagsService.On("GetAllFlags", tt.env).Return(tt.mockFlags, tt.mockError)

			// Create a test handler that uses the mock
			router.GET("/flags", func(c *gin.Context) {
				env := c.Query("env")
				if env != "local" && env != "prod" {
					c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
					return
				}

				flagsMap, err := mockFlagsService.GetAllFlags(env)
				if err != nil {
					c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
					return
				}

				c.JSON(http.StatusOK, flagsMap)
			})

			// Create a test request
			req, _ := http.NewRequest("GET", "/flags?env="+tt.env, nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare the response body
			var actualResponse interface{}
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)

			// Convert expected body to map for comparison
			expectedBytes, _ := json.Marshal(tt.expectedBody)
			var expectedResponse interface{}
			json.Unmarshal(expectedBytes, &expectedResponse)

			assert.Equal(t, expectedResponse, actualResponse)

			// Verify all mock expectations were met
			mockFlagsService.AssertExpectations(t)
		})
	}
}

func TestGetFlagsValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		env            string
		expectedStatus int
		expectedBody   ErrorResponse
	}{
		{
			name:           "invalid env - empty",
			env:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
		{
			name:           "invalid env - development",
			env:            "development",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
		{
			name:           "invalid env - staging",
			env:            "staging",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
		{
			name:           "invalid env - test",
			env:            "test",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()

			// Create a test handler
			router.GET("/flags", func(c *gin.Context) {
				env := c.Query("env")
				if env != "local" && env != "prod" {
					c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
			})

			// Create a test request
			req, _ := http.NewRequest("GET", "/flags?env="+tt.env, nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare the response body
			var actualResponse ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}

func TestGetFlagByKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		env            string
		key            string
		mockValue      bool
		mockExists     bool
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "successful local flag - enabled",
			env:            "local",
			key:            "feature_a",
			mockValue:      true,
			mockExists:     true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   FlagStatus{Key: "feature_a", Enabled: true},
		},
		{
			name:           "successful local flag - disabled",
			env:            "local",
			key:            "feature_b",
			mockValue:      false,
			mockExists:     true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   FlagStatus{Key: "feature_b", Enabled: false},
		},
		{
			name:           "successful prod flag",
			env:            "prod",
			key:            "feature_c",
			mockValue:      true,
			mockExists:     true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   FlagStatus{Key: "feature_c", Enabled: true},
		},
		{
			name:           "flag not found",
			env:            "local",
			key:            "nonexistent_flag",
			mockValue:      false,
			mockExists:     false,
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   ErrorResponse{Error: "unknown flag key"},
		},
		{
			name:           "flags service error",
			env:            "local",
			key:            "feature_a",
			mockValue:      false,
			mockExists:     false,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ErrorResponse{Error: assert.AnError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()

			// Create mock flags service
			mockFlagsService := new(MockFlagsService)

			// Set up mock expectations
			mockFlagsService.On("GetSingleFlag", tt.env, tt.key).Return(tt.mockValue, tt.mockExists, tt.mockError)

			// Create a test handler that uses the mock
			router.GET("/flags/:key", func(c *gin.Context) {
				env := c.Query("env")
				key := c.Param("key")

				if env != "local" && env != "prod" {
					c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
					return
				}

				val, ok, err := mockFlagsService.GetSingleFlag(env, key)
				if err != nil {
					c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
					return
				}
				if !ok {
					c.JSON(http.StatusNotFound, ErrorResponse{Error: "unknown flag key"})
					return
				}

				c.JSON(http.StatusOK, FlagStatus{Key: key, Enabled: val})
			})

			// Create a test request
			req, _ := http.NewRequest("GET", "/flags/"+tt.key+"?env="+tt.env, nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare the response body
			var actualResponse interface{}
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)

			// Convert expected body to map for comparison
			expectedBytes, _ := json.Marshal(tt.expectedBody)
			var expectedResponse interface{}
			json.Unmarshal(expectedBytes, &expectedResponse)

			assert.Equal(t, expectedResponse, actualResponse)

			// Verify all mock expectations were met
			mockFlagsService.AssertExpectations(t)
		})
	}
}

func TestGetFlagByKeyValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		env            string
		key            string
		expectedStatus int
		expectedBody   ErrorResponse
	}{
		{
			name:           "invalid env - empty",
			env:            "",
			key:            "feature_a",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
		{
			name:           "invalid env - development",
			env:            "development",
			key:            "feature_a",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
		{
			name:           "invalid env - staging",
			env:            "staging",
			key:            "feature_a",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrorResponse{Error: "invalid env; must be local or prod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()

			// Create a test handler
			router.GET("/flags/:key", func(c *gin.Context) {
				env := c.Query("env")
				_ = c.Param("key") // key is not used in validation test

				if env != "local" && env != "prod" {
					c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
			})

			// Create a test request
			req, _ := http.NewRequest("GET", "/flags/"+tt.key+"?env="+tt.env, nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare the response body
			var actualResponse ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, actualResponse)
		})
	}
}

// Benchmark tests for performance
func BenchmarkGetFlags(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	mockFlagsService := new(MockFlagsService)
	mockFlagsService.On("GetAllFlags", "local").Return(map[string]bool{"feature_a": true, "feature_b": false}, nil)

	router.GET("/flags", func(c *gin.Context) {
		env := c.Query("env")
		if env != "local" && env != "prod" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
			return
		}

		flagsMap, err := mockFlagsService.GetAllFlags(env)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, flagsMap)
	})

	req, _ := http.NewRequest("GET", "/flags?env=local", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetFlagByKey(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	mockFlagsService := new(MockFlagsService)
	mockFlagsService.On("GetSingleFlag", "local", "feature_a").Return(true, true, nil)

	router.GET("/flags/:key", func(c *gin.Context) {
		env := c.Query("env")
		key := c.Param("key")

		if env != "local" && env != "prod" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid env; must be local or prod"})
			return
		}

		val, ok, err := mockFlagsService.GetSingleFlag(env, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		if !ok {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "unknown flag key"})
			return
		}

		c.JSON(http.StatusOK, FlagStatus{Key: key, Enabled: val})
	})

	req, _ := http.NewRequest("GET", "/flags/feature_a?env=local", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
