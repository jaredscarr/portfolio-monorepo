package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/models"
)

// OutboxStore handles outbox event storage operations
type OutboxStore struct {
	db *DB
}

// NewOutboxStore creates a new outbox store
func NewOutboxStore(db *DB) *OutboxStore {
	return &OutboxStore{db: db}
}

// CreateEvent creates a new outbox event
func (s *OutboxStore) CreateEvent(req *models.CreateEventRequest) (*models.Event, error) {
	id := uuid.New().String()
	now := time.Now()

	// Validate that data is valid JSON
	var data interface{}
	if err := json.Unmarshal(req.Data, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON in data field: %w", err)
	}

	// Handle metadata - validate JSON if present
	var metadata interface{}
	if len(req.Metadata) > 0 {
		if err := json.Unmarshal(req.Metadata, &metadata); err != nil {
			return nil, fmt.Errorf("invalid JSON in metadata field: %w", err)
		}
		metadata = req.Metadata // Use the raw JSON for PostgreSQL
	} else {
		metadata = nil
	}

	query := `
		INSERT INTO outbox_events (id, type, source, data, metadata, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
	`

	var event models.Event
	var metadataStr sql.NullString
	var lastErrorStr sql.NullString
	var publishedAt sql.NullTime
	var dataStr string
	err := s.db.conn.QueryRow(query, id, req.Type, req.Source, req.Data, metadata, models.StatusPending, now, now).
		Scan(&event.ID, &event.Type, &event.Source, &dataStr, &metadataStr, &event.Status, &event.RetryCount, &lastErrorStr, &event.CreatedAt, &event.UpdatedAt, &publishedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Handle nullable fields
	event.Data = json.RawMessage(dataStr)

	if metadataStr.Valid {
		event.Metadata = json.RawMessage(metadataStr.String)
	} else {
		event.Metadata = nil
	}

	if lastErrorStr.Valid {
		event.LastError = lastErrorStr.String
	} else {
		event.LastError = ""
	}

	if publishedAt.Valid {
		event.PublishedAt = &publishedAt.Time
	} else {
		event.PublishedAt = nil
	}

	return &event, nil
}

// GetEvent retrieves an event by ID
func (s *OutboxStore) GetEvent(id string) (*models.Event, error) {
	query := `
		SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		WHERE id = $1
	`

	var event models.Event
	var metadataStr sql.NullString
	var lastErrorStr sql.NullString
	var publishedAt sql.NullTime
	var dataStr string
	err := s.db.conn.QueryRow(query, id).
		Scan(&event.ID, &event.Type, &event.Source, &dataStr, &metadataStr, &event.Status, &event.RetryCount, &lastErrorStr, &event.CreatedAt, &event.UpdatedAt, &publishedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Handle nullable fields
	event.Data = json.RawMessage(dataStr)

	if metadataStr.Valid {
		event.Metadata = json.RawMessage(metadataStr.String)
	} else {
		event.Metadata = nil
	}

	if lastErrorStr.Valid {
		event.LastError = lastErrorStr.String
	} else {
		event.LastError = ""
	}

	if publishedAt.Valid {
		event.PublishedAt = &publishedAt.Time
	} else {
		event.PublishedAt = nil
	}

	return &event, nil
}

// ListEvents retrieves events with pagination and filtering
func (s *OutboxStore) ListEvents(status *models.EventStatus, page, limit int) ([]models.Event, int, error) {
	offset := (page - 1) * limit

	var whereClause string
	var args []interface{}
	argIndex := 1

	if status != nil {
		whereClause = "WHERE status = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, string(*status))
		argIndex++
	}

	countQuery := "SELECT COUNT(*) FROM outbox_events " + whereClause
	var total int
	err := s.db.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	query := `
		SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, limit, offset)

	rows, err := s.db.conn.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		var metadataStr sql.NullString
		var lastErrorStr sql.NullString
		var publishedAt sql.NullTime
		var dataStr string
		err := rows.Scan(&event.ID, &event.Type, &event.Source, &dataStr, &metadataStr, &event.Status, &event.RetryCount, &lastErrorStr, &event.CreatedAt, &event.UpdatedAt, &publishedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}

		// Handle nullable fields
		event.Data = json.RawMessage(dataStr)

		if metadataStr.Valid {
			event.Metadata = json.RawMessage(metadataStr.String)
		} else {
			event.Metadata = nil
		}

		if lastErrorStr.Valid {
			event.LastError = lastErrorStr.String
		} else {
			event.LastError = ""
		}

		if publishedAt.Valid {
			event.PublishedAt = &publishedAt.Time
		} else {
			event.PublishedAt = nil
		}

		events = append(events, event)
	}

	return events, total, nil
}

// GetPendingEvents retrieves events ready for publishing
func (s *OutboxStore) GetPendingEvents(limit int) ([]models.Event, error) {
	query := `
		SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		WHERE status IN ('pending', 'retrying')
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := s.db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		var metadataStr sql.NullString
		var lastErrorStr sql.NullString
		var publishedAt sql.NullTime
		var dataStr string
		err := rows.Scan(&event.ID, &event.Type, &event.Source, &dataStr, &metadataStr, &event.Status, &event.RetryCount, &lastErrorStr, &event.CreatedAt, &event.UpdatedAt, &publishedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		// Handle nullable fields
		event.Data = json.RawMessage(dataStr)

		if metadataStr.Valid {
			event.Metadata = json.RawMessage(metadataStr.String)
		} else {
			event.Metadata = nil
		}

		if lastErrorStr.Valid {
			event.LastError = lastErrorStr.String
		} else {
			event.LastError = ""
		}

		if publishedAt.Valid {
			event.PublishedAt = &publishedAt.Time
		} else {
			event.PublishedAt = nil
		}

		events = append(events, event)
	}

	return events, nil
}

// UpdateEventStatus updates an event's status and related fields
func (s *OutboxStore) UpdateEventStatus(id string, status models.EventStatus, lastError string, retryCount int) error {
	now := time.Now()
	var publishedAt interface{}

	if status == models.StatusPublished {
		publishedAt = now
	}

	query := `
		UPDATE outbox_events
		SET status = $1, last_error = $2, retry_count = $3, updated_at = $4, published_at = $5
		WHERE id = $6
	`

	_, err := s.db.conn.Exec(query, status, lastError, retryCount, now, publishedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}

	return nil
}

// DeleteEvent deletes an event by ID
func (s *OutboxStore) DeleteEvent(id string) error {
	query := "DELETE FROM outbox_events WHERE id = $1"
	result, err := s.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
}

// GetStats returns event statistics
func (s *OutboxStore) GetStats() (*models.StatsResponse, error) {
	query := `
		SELECT 
			COUNT(*) as total_events,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_events,
			COUNT(CASE WHEN status = 'published' THEN 1 END) as published_events,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_events,
			COALESCE(SUM(retry_count), 0) as retry_count
		FROM outbox_events
	`

	var stats models.StatsResponse
	err := s.db.conn.QueryRow(query).Scan(
		&stats.TotalEvents,
		&stats.PendingEvents,
		&stats.PublishedEvents,
		&stats.FailedEvents,
		&stats.RetryCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil
}

// UpdateEventPublishedAt updates the published_at timestamp for an event
func (s *OutboxStore) UpdateEventPublishedAt(id string, publishedAt *time.Time) error {
	query := `
		UPDATE outbox_events 
		SET published_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := s.db.conn.Exec(query, publishedAt, now, id)
	if err != nil {
		return fmt.Errorf("failed to update published_at for event %s: %w", id, err)
	}

	return nil
}
