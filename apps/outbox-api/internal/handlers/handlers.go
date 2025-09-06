package handlers

import (
	"net/http"
	"strconv"

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

	// For now, just return a placeholder response
	// The actual publishing logic will be implemented in the publisher service
	c.JSON(http.StatusOK, gin.H{
		"message": "publish endpoint - implementation pending",
		"request": req,
	})
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
