package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/config"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/models"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/storage"
)

// Handler handles HTTP requests for the outbox API
type Handler struct {
	store storage.OutboxStoreInterface
	cfg   *config.Config
}

// New creates a new handler instance
func New(store storage.OutboxStoreInterface, cfg *config.Config) *Handler {
	return &Handler{
		store: store,
		cfg:   cfg,
	}
}

// CreateEvent creates a new outbox event
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

// GetEvent retrieves an event by ID
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

// RetryEvent retries publishing a failed event
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

	// Reset status to retrying
	err = h.store.UpdateEventStatus(id, models.StatusRetrying, "", event.RetryCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event queued for retry"})
}

// DeleteEvent deletes an event by ID
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

// PublishEvents publishes pending events
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
	var errors []string

	for _, event := range events {
		err := h.publishEvent(&event)
		if err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("Event %s: %v", event.ID, err))

			// Update event status to failed
			h.store.UpdateEventStatus(event.ID, models.StatusFailed, err.Error(), event.RetryCount+1)
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
		Errors:    errors,
	}

	c.JSON(http.StatusOK, response)
}

// GetStats returns service statistics
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.store.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// getEventsByIDs retrieves specific events by their IDs
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

// publishEvent sends an event to the webhook URL
func (h *Handler) publishEvent(event *models.Event) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second, // Default timeout, could be made configurable
	}

	// Prepare the webhook payload
	payload := map[string]interface{}{
		"id":         event.ID,
		"type":       event.Type,
		"source":     event.Source,
		"data":       event.Data,
		"metadata":   event.Metadata,
		"created_at": event.CreatedAt,
	}

	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", h.cfg.Publish.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "outbox-api/1.0")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
