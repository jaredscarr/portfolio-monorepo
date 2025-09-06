package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventStatus_String(t *testing.T) {
	tests := []struct {
		status   EventStatus
		expected string
	}{
		{StatusPending, "pending"},
		{StatusPublished, "published"},
		{StatusFailed, "failed"},
		{StatusRetrying, "retrying"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestCreateEventRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected CreateEventRequest
	}{
		{
			name: "valid request with metadata",
			json: `{
				"type": "user.created",
				"source": "user-service",
				"data": {"user_id": "123", "email": "test@example.com"},
				"metadata": {"version": "1.0", "correlation_id": "abc-123"}
			}`,
			expected: CreateEventRequest{
				Type:     "user.created",
				Source:   "user-service",
				Data:     json.RawMessage(`{"user_id": "123", "email": "test@example.com"}`),
				Metadata: json.RawMessage(`{"version": "1.0", "correlation_id": "abc-123"}`),
			},
		},
		{
			name: "valid request without metadata",
			json: `{
				"type": "order.placed",
				"source": "order-service",
				"data": {"order_id": "456", "amount": 99.99}
			}`,
			expected: CreateEventRequest{
				Type:     "order.placed",
				Source:   "order-service",
				Data:     json.RawMessage(`{"order_id": "456", "amount": 99.99}`),
				Metadata: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req CreateEventRequest
			err := json.Unmarshal([]byte(tt.json), &req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, req)
		})
	}
}

func TestEvent_JSON(t *testing.T) {
	now := time.Now()
	event := Event{
		ID:          "test-id-123",
		Type:        "test.event",
		Source:      "test-service",
		Data:        json.RawMessage(`{"message": "hello"}`),
		Metadata:    json.RawMessage(`{"version": "1.0"}`),
		Status:      StatusPending,
		RetryCount:  0,
		LastError:   "",
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: nil,
	}

	// Test marshaling
	jsonData, err := json.Marshal(event)
	require.NoError(t, err)

	// Test unmarshaling
	var unmarshaled Event
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, event.ID, unmarshaled.ID)
	assert.Equal(t, event.Type, unmarshaled.Type)
	assert.Equal(t, event.Source, unmarshaled.Source)
	// Compare JSON content rather than exact bytes (JSON marshaling can vary whitespace)
	assert.JSONEq(t, string(event.Data), string(unmarshaled.Data))
	assert.JSONEq(t, string(event.Metadata), string(unmarshaled.Metadata))
	assert.Equal(t, event.Status, unmarshaled.Status)
	assert.Equal(t, event.RetryCount, unmarshaled.RetryCount)
	assert.Equal(t, event.LastError, unmarshaled.LastError)
	assert.Equal(t, event.CreatedAt.Unix(), unmarshaled.CreatedAt.Unix())
	assert.Equal(t, event.UpdatedAt.Unix(), unmarshaled.UpdatedAt.Unix())
	assert.Equal(t, event.PublishedAt, unmarshaled.PublishedAt)
}

func TestEventsResponse_JSON(t *testing.T) {
	now := time.Now()
	events := []Event{
		{
			ID:        "event-1",
			Type:      "test.event",
			Source:    "test-service",
			Data:      json.RawMessage(`{"id": 1}`),
			Status:    StatusPending,
			CreatedAt: now,
		},
		{
			ID:        "event-2",
			Type:      "test.event",
			Source:    "test-service",
			Data:      json.RawMessage(`{"id": 2}`),
			Status:    StatusPublished,
			CreatedAt: now.Add(-time.Hour),
		},
	}

	response := EventsResponse{
		Events: events,
		Total:  2,
		Page:   1,
		Limit:  20,
	}

	// Test marshaling
	jsonData, err := json.Marshal(response)
	require.NoError(t, err)

	// Test unmarshaling
	var unmarshaled EventsResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, response.Total, unmarshaled.Total)
	assert.Equal(t, response.Page, unmarshaled.Page)
	assert.Equal(t, response.Limit, unmarshaled.Limit)
	assert.Len(t, unmarshaled.Events, 2)
}

func TestPublishRequest_JSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected PublishRequest
	}{
		{
			name: "publish specific events",
			json: `{
				"event_ids": ["event-1", "event-2", "event-3"],
				"batch_size": 5
			}`,
			expected: PublishRequest{
				EventIDs:  []string{"event-1", "event-2", "event-3"},
				BatchSize: 5,
			},
		},
		{
			name: "publish with batch size only",
			json: `{
				"batch_size": 10
			}`,
			expected: PublishRequest{
				EventIDs:  nil,
				BatchSize: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req PublishRequest
			err := json.Unmarshal([]byte(tt.json), &req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, req)
		})
	}
}

func TestStatsResponse_JSON(t *testing.T) {
	stats := StatsResponse{
		TotalEvents:     100,
		PendingEvents:   25,
		PublishedEvents: 70,
		FailedEvents:    5,
		RetryCount:      15,
	}

	// Test marshaling
	jsonData, err := json.Marshal(stats)
	require.NoError(t, err)

	// Test unmarshaling
	var unmarshaled StatsResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, stats, unmarshaled)
}
