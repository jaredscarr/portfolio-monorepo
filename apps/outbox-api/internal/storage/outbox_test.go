package storage

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	return &DB{conn: db}, mock
}

func TestOutboxStore_CreateEvent(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreateEventRequest
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
		expectedEvent *models.Event
	}{
		{
			name: "successful event creation",
			request: &models.CreateEventRequest{
				Type:     "test.event",
				Source:   "test-service",
				Data:     json.RawMessage(`{"message": "hello"}`),
				Metadata: json.RawMessage(`{"version": "1.0"}`),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO outbox_events (id, type, source, data, metadata, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at`).
					WithArgs(sqlmock.AnyArg(), "test.event", "test-service", sqlmock.AnyArg(), sqlmock.AnyArg(), "pending", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "source", "data", "metadata", "status", "retry_count", "last_error", "created_at", "updated_at", "published_at"}).
						AddRow("test-id", "test.event", "test-service", `{"message": "hello"}`, `{"version": "1.0"}`, "pending", 0, nil, time.Now(), time.Now(), nil))
			},
			expectedEvent: &models.Event{
				ID:          "test-id",
				Type:        "test.event",
				Source:      "test-service",
				Data:        json.RawMessage(`{"message": "hello"}`),
				Metadata:    json.RawMessage(`{"version": "1.0"}`),
				Status:      models.StatusPending,
				RetryCount:  0,
				LastError:   "",
				PublishedAt: nil,
			},
		},
		{
			name: "event creation without metadata",
			request: &models.CreateEventRequest{
				Type:     "test.event",
				Source:   "test-service",
				Data:     json.RawMessage(`{"message": "hello"}`),
				Metadata: nil,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO outbox_events (id, type, source, data, metadata, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at`).
					WithArgs(sqlmock.AnyArg(), "test.event", "test-service", sqlmock.AnyArg(), nil, "pending", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "source", "data", "metadata", "status", "retry_count", "last_error", "created_at", "updated_at", "published_at"}).
						AddRow("test-id", "test.event", "test-service", `{"message": "hello"}`, nil, "pending", 0, nil, time.Now(), time.Now(), nil))
			},
			expectedEvent: &models.Event{
				ID:          "test-id",
				Type:        "test.event",
				Source:      "test-service",
				Data:        json.RawMessage(`{"message": "hello"}`),
				Metadata:    nil,
				Status:      models.StatusPending,
				RetryCount:  0,
				LastError:   "",
				PublishedAt: nil,
			},
		},
		{
			name: "invalid JSON in data field",
			request: &models.CreateEventRequest{
				Type:     "test.event",
				Source:   "test-service",
				Data:     json.RawMessage(`invalid json`),
				Metadata: nil,
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedError: "invalid JSON in data field",
		},
		{
			name: "invalid JSON in metadata field",
			request: &models.CreateEventRequest{
				Type:     "test.event",
				Source:   "test-service",
				Data:     json.RawMessage(`{"valid": true}`),
				Metadata: json.RawMessage(`invalid json`),
			},
			mockSetup:     func(mock sqlmock.Sqlmock) {},
			expectedError: "invalid JSON in metadata field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			store := NewOutboxStore(db)
			tt.mockSetup(mock)

			event, err := store.CreateEvent(tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.expectedEvent.ID, event.ID)
				assert.Equal(t, tt.expectedEvent.Type, event.Type)
				assert.Equal(t, tt.expectedEvent.Source, event.Source)
				assert.Equal(t, tt.expectedEvent.Data, event.Data)
				assert.Equal(t, tt.expectedEvent.Metadata, event.Metadata)
				assert.Equal(t, tt.expectedEvent.Status, event.Status)
				assert.Equal(t, tt.expectedEvent.RetryCount, event.RetryCount)
				assert.Equal(t, tt.expectedEvent.LastError, event.LastError)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOutboxStore_GetEvent(t *testing.T) {
	tests := []struct {
		name          string
		eventID       string
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
		expectedEvent *models.Event
	}{
		{
			name:    "successful event retrieval",
			eventID: "test-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		WHERE id = $1`).
					WithArgs("test-id").
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "source", "data", "metadata", "status", "retry_count", "last_error", "created_at", "updated_at", "published_at"}).
						AddRow("test-id", "test.event", "test-service", `{"message": "hello"}`, `{"version": "1.0"}`, "pending", 0, nil, time.Now(), time.Now(), nil))
			},
			expectedEvent: &models.Event{
				ID:          "test-id",
				Type:        "test.event",
				Source:      "test-service",
				Data:        json.RawMessage(`{"message": "hello"}`),
				Metadata:    json.RawMessage(`{"version": "1.0"}`),
				Status:      models.StatusPending,
				RetryCount:  0,
				LastError:   "",
				PublishedAt: nil,
			},
		},
		{
			name:    "event not found",
			eventID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		WHERE id = $1`).
					WithArgs("non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "event not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			store := NewOutboxStore(db)
			tt.mockSetup(mock)

			event, err := store.GetEvent(tt.eventID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.expectedEvent.ID, event.ID)
				assert.Equal(t, tt.expectedEvent.Type, event.Type)
				assert.Equal(t, tt.expectedEvent.Source, event.Source)
				assert.Equal(t, tt.expectedEvent.Data, event.Data)
				assert.Equal(t, tt.expectedEvent.Metadata, event.Metadata)
				assert.Equal(t, tt.expectedEvent.Status, event.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOutboxStore_ListEvents(t *testing.T) {
	tests := []struct {
		name          string
		status        *models.EventStatus
		page          int
		limit         int
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
		expectedCount int
		expectedTotal int
	}{
		{
			name:   "list all events",
			status: nil,
			page:   1,
			limit:  10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Count query
				mock.ExpectQuery("SELECT COUNT(*) FROM outbox_events").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// List query
				mock.ExpectQuery(`SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`).
					WithArgs(10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "source", "data", "metadata", "status", "retry_count", "last_error", "created_at", "updated_at", "published_at"}).
						AddRow("event-1", "test.event", "test-service", `{"id": 1}`, nil, "pending", 0, nil, time.Now(), time.Now(), nil).
						AddRow("event-2", "test.event", "test-service", `{"id": 2}`, nil, "published", 0, nil, time.Now().Add(-time.Hour), time.Now().Add(-time.Hour), time.Now().Add(-time.Hour)))
			},
			expectedCount: 2,
			expectedTotal: 2,
		},
		{
			name:   "list events by status",
			status: func() *models.EventStatus { s := models.StatusPending; return &s }(),
			page:   1,
			limit:  10,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Count query
				mock.ExpectQuery("SELECT COUNT(*) FROM outbox_events WHERE status = $1").
					WithArgs("pending").
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				// List query
				mock.ExpectQuery(`SELECT id, type, source, data, metadata, status, retry_count, last_error, created_at, updated_at, published_at
		FROM outbox_events
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`).
					WithArgs("pending", 10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type", "source", "data", "metadata", "status", "retry_count", "last_error", "created_at", "updated_at", "published_at"}).
						AddRow("event-1", "test.event", "test-service", `{"id": 1}`, nil, "pending", 0, nil, time.Now(), time.Now(), nil))
			},
			expectedCount: 1,
			expectedTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			store := NewOutboxStore(db)
			tt.mockSetup(mock)

			events, total, err := store.ListEvents(tt.status, tt.page, tt.limit)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Len(t, events, tt.expectedCount)
				assert.Equal(t, tt.expectedTotal, total)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOutboxStore_UpdateEventStatus(t *testing.T) {
	tests := []struct {
		name          string
		eventID       string
		status        models.EventStatus
		lastError     string
		retryCount    int
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
	}{
		{
			name:       "successful status update",
			eventID:    "test-id",
			status:     models.StatusPublished,
			lastError:  "",
			retryCount: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE outbox_events
		SET status = $1, last_error = $2, retry_count = $3, updated_at = $4, published_at = $5
		WHERE id = $6`).
					WithArgs("published", "", 0, sqlmock.AnyArg(), sqlmock.AnyArg(), "test-id").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:       "update to failed status",
			eventID:    "test-id",
			status:     models.StatusFailed,
			lastError:  "connection timeout",
			retryCount: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE outbox_events
		SET status = $1, last_error = $2, retry_count = $3, updated_at = $4, published_at = $5
		WHERE id = $6`).
					WithArgs("failed", "connection timeout", 3, sqlmock.AnyArg(), nil, "test-id").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			store := NewOutboxStore(db)
			tt.mockSetup(mock)

			err := store.UpdateEventStatus(tt.eventID, tt.status, tt.lastError, tt.retryCount)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOutboxStore_DeleteEvent(t *testing.T) {
	tests := []struct {
		name          string
		eventID       string
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
	}{
		{
			name:    "successful deletion",
			eventID: "test-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM outbox_events WHERE id = $1").
					WithArgs("test-id").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:    "event not found",
			eventID: "non-existent-id",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM outbox_events WHERE id = $1").
					WithArgs("non-existent-id").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: "event not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			store := NewOutboxStore(db)
			tt.mockSetup(mock)

			err := store.DeleteEvent(tt.eventID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestOutboxStore_GetStats(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedError string
		expectedStats *models.StatsResponse
	}{
		{
			name: "successful stats retrieval",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT 
			COUNT(*) as total_events,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_events,
			COUNT(CASE WHEN status = 'published' THEN 1 END) as published_events,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_events,
			SUM(retry_count) as retry_count
		FROM outbox_events`).
					WillReturnRows(sqlmock.NewRows([]string{"total_events", "pending_events", "published_events", "failed_events", "retry_count"}).
						AddRow(100, 25, 70, 5, 15))
			},
			expectedStats: &models.StatsResponse{
				TotalEvents:     100,
				PendingEvents:   25,
				PublishedEvents: 70,
				FailedEvents:    5,
				RetryCount:      15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			defer db.Close()

			store := NewOutboxStore(db)
			tt.mockSetup(mock)

			stats, err := store.GetStats()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, stats)
				assert.Equal(t, tt.expectedStats, stats)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
