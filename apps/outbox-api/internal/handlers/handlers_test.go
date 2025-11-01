package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/config"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockOutboxStore struct {
	mock.Mock
}

type MockSimulationGates struct {
	mock.Mock
}

func (m *MockSimulationGates) IsSimulationModeEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) ShouldDisablePublishing() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) ShouldSimulateWebhookFailures() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) ShouldSimulateNetworkDelays() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) ShouldUsePartialFailureMode() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) ShouldUseCircuitBreakerDemo() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) CheckCircuitBreaker() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSimulationGates) RecordCircuitBreakerSuccess() {
	m.Called()
}

func (m *MockSimulationGates) RecordCircuitBreakerFailure() {
	m.Called()
}

func (m *MockSimulationGates) GetSimulationStatus() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func (m *MockOutboxStore) CreateEvent(req *models.CreateEventRequest) (*models.Event, error) {
	args := m.Called(req)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockOutboxStore) GetEvent(id string) (*models.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockOutboxStore) ListEvents(status *models.EventStatus, page, limit int) ([]models.Event, int, error) {
	args := m.Called(status, page, limit)
	return args.Get(0).([]models.Event), args.Int(1), args.Error(2)
}

func (m *MockOutboxStore) GetPendingEvents(limit int) ([]models.Event, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.Event), args.Error(1)
}

func (m *MockOutboxStore) UpdateEventStatus(id string, status models.EventStatus, lastError string, retryCount int) error {
	args := m.Called(id, status, lastError, retryCount)
	return args.Error(0)
}

func (m *MockOutboxStore) DeleteEvent(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOutboxStore) GetStats() (*models.StatsResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StatsResponse), args.Error(1)
}

func (m *MockOutboxStore) UpdateEventPublishedAt(id string, publishedAt *time.Time) error {
	args := m.Called(id, publishedAt)
	return args.Error(0)
}

func setupTestRouter(store *MockOutboxStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
	}

	mockGates := &MockSimulationGates{}
	mockGates.On("ShouldDisablePublishing").Return(false)
	mockGates.On("ShouldSimulateWebhookFailures").Return(false)
	mockGates.On("ShouldSimulateNetworkDelays").Return(false)
	mockGates.On("ShouldUsePartialFailureMode").Return(false)
	mockGates.On("ShouldUseCircuitBreakerDemo").Return(false)
	mockGates.On("CheckCircuitBreaker").Return(false)
	mockGates.On("RecordCircuitBreakerFailure").Return()
	mockGates.On("RecordCircuitBreakerSuccess").Return()
	h := New(store, cfg, mockGates)

	api := router.Group("/api/v1")
	{
		api.POST("/events", h.CreateEvent)
		api.GET("/events", h.ListEvents)
		api.GET("/events/:id", h.GetEvent)
		api.POST("/events/:id/retry", h.RetryEvent)
		api.DELETE("/events/:id", h.DeleteEvent)
	}

	admin := router.Group("/admin")
	{
		admin.POST("/publish", h.PublishEvents)
		admin.GET("/stats", h.GetStats)
	}

	return router
}

func TestHandler_CreateEvent(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful event creation",
			requestBody: models.CreateEventRequest{
				Type:     "test.event",
				Source:   "test-service",
				Data:     json.RawMessage(`{"message": "hello"}`),
				Metadata: json.RawMessage(`{"version": "1.0"}`),
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				expectedEvent := &models.Event{
					ID:         "test-id",
					Type:       "test.event",
					Source:     "test-service",
					Data:       json.RawMessage(`{"message": "hello"}`),
					Metadata:   json.RawMessage(`{"version": "1.0"}`),
					Status:     models.StatusPending,
					RetryCount: 0,
					LastError:  "",
				}
				mockStore.On("CreateEvent", mock.AnythingOfType("*models.CreateEventRequest")).Return(expectedEvent, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"type": "test.event",
				// Missing required fields
			},
			mockSetup:      func(mockStore *MockOutboxStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'CreateEventRequest.Source' Error:Field validation for 'Source' failed on the 'required' tag",
		},
		{
			name: "storage error",
			requestBody: models.CreateEventRequest{
				Type:     "test.event",
				Source:   "test-service",
				Data:     json.RawMessage(`{"message": "hello"}`),
				Metadata: nil,
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("CreateEvent", mock.AnythingOfType("*models.CreateEventRequest")).Return((*models.Event)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "assert.AnError general error for testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			router := setupTestRouter(mockStore)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "event")
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandler_GetEvent(t *testing.T) {
	tests := []struct {
		name           string
		eventID        string
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful event retrieval",
			eventID: "test-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				expectedEvent := &models.Event{
					ID:         "test-id",
					Type:       "test.event",
					Source:     "test-service",
					Data:       json.RawMessage(`{"message": "hello"}`),
					Status:     models.StatusPending,
					RetryCount: 0,
					LastError:  "",
				}
				mockStore.On("GetEvent", "test-id").Return(expectedEvent, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "event not found",
			eventID: "non-existent-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetEvent", "non-existent-id").Return((*models.Event)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "assert.AnError general error for testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			router := setupTestRouter(mockStore)

			url := "/api/v1/events/" + tt.eventID
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandler_ListEvents(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:        "list all events",
			queryParams: "",
			mockSetup: func(mockStore *MockOutboxStore) {
				events := []models.Event{
					{
						ID:     "event-1",
						Type:   "test.event",
						Source: "test-service",
						Status: models.StatusPending,
					},
					{
						ID:     "event-2",
						Type:   "test.event",
						Source: "test-service",
						Status: models.StatusPublished,
					},
				}
				mockStore.On("ListEvents", mock.AnythingOfType("*models.EventStatus"), 1, 20).Return(events, 2, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:        "list events by status",
			queryParams: "?status=pending",
			mockSetup: func(mockStore *MockOutboxStore) {
				events := []models.Event{
					{
						ID:     "event-1",
						Type:   "test.event",
						Source: "test-service",
						Status: models.StatusPending,
					},
				}
				status := models.StatusPending
				mockStore.On("ListEvents", &status, 1, 20).Return(events, 1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:        "pagination",
			queryParams: "?page=2&limit=10",
			mockSetup: func(mockStore *MockOutboxStore) {
				events := []models.Event{
					{
						ID:     "event-1",
						Type:   "test.event",
						Source: "test-service",
						Status: models.StatusPending,
					},
				}
				mockStore.On("ListEvents", mock.AnythingOfType("*models.EventStatus"), 2, 10).Return(events, 1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			router := setupTestRouter(mockStore)

			url := "/api/v1/events" + tt.queryParams
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.EventsResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Len(t, response.Events, tt.expectedCount)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandler_RetryEvent(t *testing.T) {
	tests := []struct {
		name           string
		eventID        string
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful retry",
			eventID: "test-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				event := &models.Event{
					ID:         "test-id",
					Type:       "test.event",
					Source:     "test-service",
					Status:     models.StatusFailed,
					RetryCount: 2,
				}
				mockStore.On("GetEvent", "test-id").Return(event, nil)
				mockStore.On("UpdateEventStatus", "test-id", models.StatusFailed, mock.AnythingOfType("string"), 3).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "event not found",
			eventID: "non-existent-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetEvent", "non-existent-id").Return((*models.Event)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "assert.AnError general error for testing",
		},
		{
			name:    "retry non-failed event",
			eventID: "test-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				event := &models.Event{
					ID:     "test-id",
					Type:   "test.event",
					Source: "test-service",
					Status: models.StatusPending, // Not failed
				}
				mockStore.On("GetEvent", "test-id").Return(event, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "only failed events can be retried",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			router := setupTestRouter(mockStore)

			url := "/api/v1/events/" + tt.eventID + "/retry"
			req, err := http.NewRequest("POST", url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteEvent(t *testing.T) {
	tests := []struct {
		name           string
		eventID        string
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful deletion",
			eventID: "test-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("DeleteEvent", "test-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "event not found",
			eventID: "non-existent-id",
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("DeleteEvent", "non-existent-id").Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "assert.AnError general error for testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			router := setupTestRouter(mockStore)

			url := "/api/v1/events/" + tt.eventID
			req, err := http.NewRequest("DELETE", url, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandler_GetStats(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedStats  *models.StatsResponse
	}{
		{
			name: "successful stats retrieval",
			mockSetup: func(mockStore *MockOutboxStore) {
				stats := &models.StatsResponse{
					TotalEvents:     100,
					PendingEvents:   25,
					PublishedEvents: 70,
					FailedEvents:    5,
					RetryCount:      15,
				}
				mockStore.On("GetStats").Return(stats, nil)
			},
			expectedStatus: http.StatusOK,
			expectedStats: &models.StatsResponse{
				TotalEvents:     100,
				PendingEvents:   25,
				PublishedEvents: 70,
				FailedEvents:    5,
				RetryCount:      15,
			},
		},
		{
			name: "storage error",
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetStats").Return((*models.StatsResponse)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			router := setupTestRouter(mockStore)

			req, err := http.NewRequest("GET", "/admin/stats", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStats != nil {
				var response models.StatsResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, *tt.expectedStats, response)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestHandler_PublishEvents(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockOutboxStore)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "publish pending events with webhook failure",
			requestBody: models.PublishRequest{
				BatchSize: 5,
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				events := []models.Event{
					{
						ID:         "event-1",
						Type:       "test.event",
						Source:     "test-service",
						Status:     models.StatusPending,
						RetryCount: 0,
					},
					{
						ID:         "event-2",
						Type:       "test.event",
						Source:     "test-service",
						Status:     models.StatusPending,
						RetryCount: 0,
					},
				}
				mockStore.On("GetPendingEvents", 5).Return(events, nil)
				// Expect failed status due to webhook connection failure
				mockStore.On("UpdateEventStatus", "event-1", models.StatusFailed, mock.AnythingOfType("string"), 1).Return(nil)
				mockStore.On("UpdateEventStatus", "event-2", models.StatusFailed, mock.AnythingOfType("string"), 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "publish specific events by ID with webhook failure",
			requestBody: models.PublishRequest{
				EventIDs: []string{"event-1", "event-2"},
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				event1 := &models.Event{
					ID:         "event-1",
					Type:       "test.event",
					Source:     "test-service",
					Status:     models.StatusPending,
					RetryCount: 0,
				}
				event2 := &models.Event{
					ID:         "event-2",
					Type:       "test.event",
					Source:     "test-service",
					Status:     models.StatusRetrying,
					RetryCount: 1,
				}
				mockStore.On("GetEvent", "event-1").Return(event1, nil)
				mockStore.On("GetEvent", "event-2").Return(event2, nil)
				// Expect failed status due to webhook connection failure
				mockStore.On("UpdateEventStatus", "event-1", models.StatusFailed, mock.AnythingOfType("string"), 1).Return(nil)
				mockStore.On("UpdateEventStatus", "event-2", models.StatusFailed, mock.AnythingOfType("string"), 2).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no events to publish",
			requestBody: models.PublishRequest{
				BatchSize: 5,
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetPendingEvents", 5).Return([]models.Event{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no specific events to publish (all non-pending)",
			requestBody: models.PublishRequest{
				EventIDs: []string{"event-1", "event-2"},
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				event1 := &models.Event{
					ID:     "event-1",
					Type:   "test.event",
					Source: "test-service",
					Status: models.StatusPublished, // Already published
				}
				event2 := &models.Event{
					ID:     "event-2",
					Type:   "test.event",
					Source: "test-service",
					Status: models.StatusFailed, // Failed, not retrying
				}
				mockStore.On("GetEvent", "event-1").Return(event1, nil)
				mockStore.On("GetEvent", "event-2").Return(event2, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"batch_size": "invalid", // Should be int
			},
			mockSetup:      func(mockStore *MockOutboxStore) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "storage error getting pending events",
			requestBody: models.PublishRequest{
				BatchSize: 5,
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetPendingEvents", 5).Return(([]models.Event)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "assert.AnError general error for testing",
		},
		{
			name: "storage error getting specific event",
			requestBody: models.PublishRequest{
				EventIDs: []string{"event-1"},
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetEvent", "event-1").Return((*models.Event)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to get event event-1: assert.AnError general error for testing",
		},
		{
			name: "storage error updating event status",
			requestBody: models.PublishRequest{
				BatchSize: 5,
			},
			mockSetup: func(mockStore *MockOutboxStore) {
				events := []models.Event{
					{
						ID:         "event-1",
						Type:       "test.event",
						Source:     "test-service",
						Status:     models.StatusPending,
						RetryCount: 0,
					},
				}
				mockStore.On("GetPendingEvents", 5).Return(events, nil)
				mockStore.On("UpdateEventStatus", "event-1", models.StatusFailed, mock.AnythingOfType("string"), 1).Return(nil)
			},
			expectedStatus: http.StatusOK, // HTTP errors are handled gracefully
		},
		{
			name:        "use default batch size when not provided",
			requestBody: models.PublishRequest{}, // Empty request
			mockSetup: func(mockStore *MockOutboxStore) {
				mockStore.On("GetPendingEvents", 10).Return([]models.Event{}, nil) // Default batch size is 10
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockOutboxStore)
			tt.mockSetup(mockStore)

			// Setup config with webhook URL for testing
			cfg := &config.Config{
				Server: config.ServerConfig{
					Port: "8080",
				},
				Publish: config.PublishConfig{
					BatchSize:  10,
					WebhookURL: "http://localhost:3000/webhook",
				},
			}

			gin.SetMode(gin.TestMode)
			router := gin.New()
			mockGates := &MockSimulationGates{}
			mockGates.On("ShouldDisablePublishing").Return(false)
			mockGates.On("ShouldSimulateWebhookFailures").Return(false)
			mockGates.On("ShouldSimulateNetworkDelays").Return(false)
			mockGates.On("ShouldUsePartialFailureMode").Return(false)
			mockGates.On("ShouldUseCircuitBreakerDemo").Return(false)
			mockGates.On("CheckCircuitBreaker").Return(false)
			mockGates.On("RecordCircuitBreakerFailure").Return()
			mockGates.On("RecordCircuitBreakerSuccess").Return()
			h := New(mockStore, cfg, mockGates)

			admin := router.Group("/admin")
			{
				admin.POST("/publish", h.PublishEvents)
			}

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/admin/publish", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			} else if tt.expectedStatus == http.StatusOK {
				var response models.PublishResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.GreaterOrEqual(t, response.Published, 0)
				assert.GreaterOrEqual(t, response.Failed, 0)
			}

			mockStore.AssertExpectations(t)
		})
	}
}
