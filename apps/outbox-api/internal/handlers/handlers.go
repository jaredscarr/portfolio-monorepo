package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/config"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/gates"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/models"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/storage"
)

// Special error to indicate publishing was skipped due to simulation
var ErrPublishingSkipped = errors.New("publishing skipped due to simulation")

// Handler handles HTTP requests for the outbox API
type Handler struct {
	store           storage.OutboxStoreInterface
	cfg             *config.Config
	simulationGates gates.SimulationGatesInterface
}

// New creates a new handler instance
func New(store storage.OutboxStoreInterface, cfg *config.Config, simulationGates gates.SimulationGatesInterface) *Handler {
	return &Handler{
		store:           store,
		cfg:             cfg,
		simulationGates: simulationGates,
	}
}

// CreateEvent godoc
// @Summary Create a new outbox event
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/events [post]
func (h *Handler) CreateEvent(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.store.CreateEvent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}

// GetEvent godoc
// @Summary Get event by ID
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/events/{id} [get]
func (h *Handler) GetEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	event, err := h.store.GetEvent(id)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// ListEvents retrieves events with pagination and filtering
func (h *Handler) ListEvents(c *gin.Context) {
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	var status *models.EventStatus
	if s := c.Query("status"); s != "" {
		eventStatus := models.EventStatus(s)
		if eventStatus == models.StatusPending || eventStatus == models.StatusPublished || eventStatus == models.StatusFailed || eventStatus == models.StatusRetrying {
			status = &eventStatus
		}
	}

	events, total, err := h.store.ListEvents(status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.EventsResponse{
		Events: events,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RetryEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	event, err := h.store.GetEvent(id)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if event.Status != models.StatusFailed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only failed events can be retried"})
		return
	}

	// Immediately attempt to publish the event
	err = h.publishEvent(event)
	if err != nil {
		// Publishing failed - update to failed status with incremented retry count
		h.store.UpdateEventStatus(id, models.StatusFailed, err.Error(), event.RetryCount+1)
		c.JSON(http.StatusOK, gin.H{
			"message": "retry attempted but failed", 
			"error": err.Error(),
			"retry_count": event.RetryCount + 1,
		})
		return
	}

	// Publishing succeeded - update to published status
	now := time.Now()
	h.store.UpdateEventStatus(id, models.StatusPublished, "", event.RetryCount)
	h.store.UpdateEventPublishedAt(id, &now)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "event retried and published successfully",
		"retry_count": event.RetryCount,
	})
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event ID is required"})
		return
	}

	err := h.store.DeleteEvent(id)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}

func (h *Handler) PublishEvents(c *gin.Context) {
	var req models.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default batch size if not provided
	batchSize := req.BatchSize
	if batchSize <= 0 {
		batchSize = h.cfg.Publish.BatchSize
	}

	var events []models.Event
	var err error

	// Get events to publish
	if len(req.EventIDs) > 0 {
		// Get specific events by ID
		events, err = h.getEventsByIDs(req.EventIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Get pending events
		events, err = h.store.GetPendingEvents(batchSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(events) == 0 {
		c.JSON(http.StatusOK, models.PublishResponse{
			Published: 0,
			Failed:    0,
		})
		return
	}

	// Publish events
	published := 0
	failed := 0
	var errorMessages []string

	for i, event := range events {
		var err error
		
		// Check for partial failure simulation
		if h.simulationGates.ShouldUsePartialFailureMode() {
			// Simulate partial failures - every 3rd event fails, others succeed
			if i%3 == 2 {
				err = fmt.Errorf("simulated partial batch failure (event %d in batch)", i+1)
			} else {
				// Simulate success without actual webhook call
				err = nil
				fmt.Printf("DEBUG: Simulated success for event %s (partial failure mode)\n", event.ID)
			}
		} else {
			err = h.publishEvent(&event)
		}
		
		if err != nil {
			if errors.Is(err, ErrPublishingSkipped) {
				// Publishing was skipped due to simulation - don't count as published or failed
				fmt.Printf("DEBUG: Event %s skipped due to simulation\n", event.ID)
				// Event stays in pending state - no status update needed
			} else {
				// Actual failure
				failed++
				errorMessages = append(errorMessages, fmt.Sprintf("Event %s: %v", event.ID, err))

				// Update event status to failed
				h.store.UpdateEventStatus(event.ID, models.StatusFailed, err.Error(), event.RetryCount+1)
			}
		} else {
			published++

			// Update event status to published
			now := time.Now()
			h.store.UpdateEventStatus(event.ID, models.StatusPublished, "", event.RetryCount)
			h.store.UpdateEventPublishedAt(event.ID, &now)
		}
	}

	response := models.PublishResponse{
		Published: published,
		Failed:    failed,
		Errors:    errorMessages,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.store.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetSimulationStatus(c *gin.Context) {
	status := h.simulationGates.GetSimulationStatus()
	c.JSON(http.StatusOK, gin.H{
		"simulation_status": status,
	})
}

func (h *Handler) getEventsByIDs(eventIDs []string) ([]models.Event, error) {
	var events []models.Event

	for _, id := range eventIDs {
		event, err := h.store.GetEvent(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get event %s: %w", id, err)
		}

		// Only include events that are pending or retrying
		if event.Status == models.StatusPending || event.Status == models.StatusRetrying {
			events = append(events, *event)
		}
	}

	return events, nil
}

func (h *Handler) publishEvent(event *models.Event) error {

	shouldDisable := h.simulationGates.ShouldDisablePublishing()
	fmt.Printf("DEBUG: ShouldDisablePublishing() = %v for event %s\n", shouldDisable, event.ID)
	
	if shouldDisable {
		fmt.Printf("DEBUG: Publishing disabled by simulation gate for event %s\n", event.ID)
		return ErrPublishingSkipped
	}

	if h.simulationGates.CheckCircuitBreaker() {
		// Circuit is open - fail fast without making request
		fmt.Printf("DEBUG: Circuit breaker OPEN - failing fast for event %s\n", event.ID)
		return fmt.Errorf("circuit breaker is open - request blocked")
	}

	if h.simulationGates.ShouldSimulateNetworkDelays() {
		// Add artificial delay to simulate network issues
		time.Sleep(2 * time.Second)
	}

	// Check for forced failures AFTER circuit breaker and delays
	if h.simulationGates.ShouldSimulateWebhookFailures() {
		h.simulationGates.RecordCircuitBreakerFailure() // Record failure for circuit breaker
		return fmt.Errorf("simulated webhook failure (forced by feature gate)")
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // Default timeout, could be made configurable
	}

	payload := map[string]interface{}{
		"id":         event.ID,
		"type":       event.Type,
		"source":     event.Source,
		"data":       event.Data,
		"metadata":   event.Metadata,
		"created_at": event.CreatedAt,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	req, err := http.NewRequest("POST", h.cfg.Publish.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "outbox-api/1.0")

	resp, err := client.Do(req)
	if err != nil {
		h.simulationGates.RecordCircuitBreakerFailure()
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		h.simulationGates.RecordCircuitBreakerFailure()
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	h.simulationGates.RecordCircuitBreakerSuccess()
	return nil
}
