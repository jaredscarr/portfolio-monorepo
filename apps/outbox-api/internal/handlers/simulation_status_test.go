package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetSimulationStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStore := &MockOutboxStore{}
	mockGates := &MockSimulationGates{}

	expectedStatus := map[string]interface{}{
		"simulation_mode_enabled": true,
		"disable_publishing":      false,
		"force_webhook_failures":  true,
	}

	mockGates.On("GetSimulationStatus").Return(expectedStatus)

	cfg := &config.Config{}
	h := New(mockStore, cfg, mockGates)

	router := gin.New()
	router.GET("/admin/simulation-status", h.GetSimulationStatus)

	req, _ := http.NewRequest("GET", "/admin/simulation-status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	simulationStatus := response["simulation_status"].(map[string]interface{})
	assert.Equal(t, true, simulationStatus["simulation_mode_enabled"])
	assert.Equal(t, false, simulationStatus["disable_publishing"])
	assert.Equal(t, true, simulationStatus["force_webhook_failures"])

	mockGates.AssertExpectations(t)
}
