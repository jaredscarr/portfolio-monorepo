package models

import (
	"encoding/json"
	"time"
)

// EventStatus represents the current status of an event
type EventStatus string

const (
	StatusPending   EventStatus = "pending"
	StatusPublished EventStatus = "published"
	StatusFailed    EventStatus = "failed"
	StatusRetrying  EventStatus = "retrying"
)

// Event represents an outbox event
type Event struct {
	ID          string          `json:"id" db:"id"`
	Type        string          `json:"type" db:"type"`
	Source      string          `json:"source" db:"source"`
	Data        json.RawMessage `json:"data" db:"data"`
	Metadata    json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	Status      EventStatus     `json:"status" db:"status"`
	RetryCount  int             `json:"retry_count" db:"retry_count"`
	LastError   string          `json:"last_error,omitempty" db:"last_error"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	PublishedAt *time.Time      `json:"published_at,omitempty" db:"published_at"`
}

// CreateEventRequest represents the request to create a new event
type CreateEventRequest struct {
	Type     string          `json:"type" binding:"required"`
	Source   string          `json:"source" binding:"required"`
	Data     json.RawMessage `json:"data" binding:"required"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// EventResponse represents the response for event operations
type EventResponse struct {
	Event *Event `json:"event,omitempty"`
	Error string `json:"error,omitempty"`
}

// EventsResponse represents the response for listing events
type EventsResponse struct {
	Events []Event `json:"events"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

// PublishRequest represents the request to publish events
type PublishRequest struct {
	EventIDs  []string `json:"event_ids,omitempty"`
	BatchSize int      `json:"batch_size,omitempty"`
}

// PublishResponse represents the response for publish operations
type PublishResponse struct {
	Published int      `json:"published"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
}

// StatsResponse represents service statistics
type StatsResponse struct {
	TotalEvents     int `json:"total_events"`
	PendingEvents   int `json:"pending_events"`
	PublishedEvents int `json:"published_events"`
	FailedEvents    int `json:"failed_events"`
	RetryCount      int `json:"retry_count"`
}
