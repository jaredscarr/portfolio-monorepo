package storage

import (
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/models"
	"time"
)

// OutboxStoreInterface defines the interface for outbox event storage operations
type OutboxStoreInterface interface {
	CreateEvent(req *models.CreateEventRequest) (*models.Event, error)
	GetEvent(id string) (*models.Event, error)
	ListEvents(status *models.EventStatus, page, limit int) ([]models.Event, int, error)
	GetPendingEvents(limit int) ([]models.Event, error)
	UpdateEventStatus(id string, status models.EventStatus, lastError string, retryCount int) error
	UpdateEventPublishedAt(id string, publishedAt *time.Time) error
	DeleteEvent(id string) error
	GetStats() (*models.StatsResponse, error)
}
