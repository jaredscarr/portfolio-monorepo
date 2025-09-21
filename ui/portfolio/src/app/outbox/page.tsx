"use client";

import React, { useState, useEffect, useCallback } from "react";
import {
  Container,
  Typography,
  Box,
  Button,
  Paper,
  Alert,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Pagination,
  Snackbar,
} from "@mui/material";
import {
  Add as AddIcon,
  Refresh as RefreshIcon,
  Publish as PublishIcon,
} from "@mui/icons-material";
import { Navigation } from "../../components/Navigation";
import { EventsTable } from "../../components/outbox/EventsTable";
import { CreateEventForm } from "../../components/outbox/CreateEventForm";
import { EventDetailDialog } from "../../components/outbox/EventDetailDialog";
import { StatsCards } from "../../components/outbox/StatsCards";
import { SimulationControls } from "../../components/outbox/SimulationControls";
import {
  OutboxEvent,
  EventsResponse,
  CreateEventRequest,
  OutboxStats,
  PublishRequest,
} from "../../types/outbox";

export default function OutboxPage() {
  const [events, setEvents] = useState<OutboxEvent[]>([]);
  const [stats, setStats] = useState<OutboxStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Pagination and filtering
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string>("");

  // Dialogs
  const [createFormOpen, setCreateFormOpen] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<OutboxEvent | null>(null);
  const [detailDialogOpen, setDetailDialogOpen] = useState(false);

  const fetchEvents = useCallback(async () => {
    try {
      setLoading(true);
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
      });

      if (statusFilter) {
        params.set("status", statusFilter);
      }

      const response = await fetch(`/api/outbox/events?${params}`);
      if (!response.ok) {
        throw new Error("Failed to fetch events");
      }

      const data: EventsResponse = await response.json();
      setEvents(data.events || []);
      setTotalPages(Math.ceil(data.total / limit));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch events");
    } finally {
      setLoading(false);
    }
  }, [page, limit, statusFilter]);

  const fetchStats = async () => {
    try {
      const response = await fetch("/api/outbox/admin/stats");
      if (!response.ok) {
        throw new Error("Failed to fetch stats");
      }

      const data: OutboxStats = await response.json();
      setStats(data);
    } catch (err) {
      console.error("Failed to fetch stats:", err);
    }
  };

  useEffect(() => {
    fetchEvents();
    fetchStats();
  }, [fetchEvents]);

  const handleCreateEvent = async (eventData: CreateEventRequest) => {
    try {
      const response = await fetch("/api/outbox/events", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(eventData),
      });

      if (!response.ok) {
        throw new Error("Failed to create event");
      }

      setSuccess("Event created successfully");
      fetchEvents();
      fetchStats();
    } catch (err) {
      throw err; // Let the form handle the error
    }
  };

  const handleRetryEvent = async (eventId: string) => {
    try {
      const response = await fetch(`/api/outbox/events/${eventId}/retry`, {
        method: "POST",
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to retry event");
      }

      const result = await response.json();
      setSuccess(result.message);
      fetchEvents();
      fetchStats();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to retry event");
    }
  };

  const handleDeleteEvent = async (eventId: string) => {
    if (!confirm("Are you sure you want to delete this event?")) {
      return;
    }

    try {
      const response = await fetch(`/api/outbox/events/${eventId}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error("Failed to delete event");
      }

      setSuccess("Event deleted successfully");
      fetchEvents();
      fetchStats();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete event");
    }
  };

  const handlePublishEvents = async () => {
    try {
      setLoading(true);
      const publishRequest: PublishRequest = {
        batch_size: 10,
      };

      const response = await fetch("/api/outbox/admin/publish", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(publishRequest),
      });

      if (!response.ok) {
        throw new Error("Failed to publish events");
      }

      const result = await response.json();
      setSuccess(
        `Published ${result.published} events, ${result.failed} failed`
      );
      fetchEvents();
      fetchStats();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to publish events");
    } finally {
      setLoading(false);
    }
  };

  const handleViewEvent = (event: OutboxEvent) => {
    setSelectedEvent(event);
    setDetailDialogOpen(true);
  };

  return (
    <>
      <Navigation />
      <Container maxWidth="xl" sx={{ py: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            Outbox Event Management
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Monitor and manage outbox pattern events for reliable message
            delivery
          </Typography>
        </Box>

        {error && (
          <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        <StatsCards stats={stats} />

        <Box sx={{ mt: 4 }}>
          <SimulationControls
            onSimulationChange={() => {
              fetchStats();
              fetchEvents();
            }}
          />
        </Box>

        <Paper sx={{ p: 3, mt: 4 }}>
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              mb: 3,
            }}
          >
            <Typography variant="h6">Events</Typography>
            <Box sx={{ display: "flex", gap: 2 }}>
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => setCreateFormOpen(true)}
              >
                Create Event
              </Button>
              <Button
                variant="outlined"
                startIcon={<PublishIcon />}
                onClick={handlePublishEvents}
                disabled={loading}
              >
                Publish Pending
              </Button>
              <Button
                variant="outlined"
                startIcon={<RefreshIcon />}
                onClick={() => {
                  fetchEvents();
                  fetchStats();
                }}
                disabled={loading}
              >
                Refresh
              </Button>
            </Box>
          </Box>

          <Box sx={{ display: "flex", gap: 2, mb: 3 }}>
            <FormControl size="small" sx={{ minWidth: 120 }}>
              <InputLabel>Status</InputLabel>
              <Select
                value={statusFilter}
                label="Status"
                onChange={(e) => {
                  setStatusFilter(e.target.value);
                  setPage(1);
                }}
              >
                <MenuItem value="">All</MenuItem>
                <MenuItem value="pending">Pending</MenuItem>
                <MenuItem value="published">Published</MenuItem>
                <MenuItem value="failed">Failed</MenuItem>
                <MenuItem value="retrying">Retrying</MenuItem>
              </Select>
            </FormControl>
          </Box>

          {loading && events.length === 0 ? (
            <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
              <CircularProgress />
            </Box>
          ) : (
            <EventsTable
              events={events}
              loading={loading}
              onRetry={handleRetryEvent}
              onDelete={handleDeleteEvent}
              onView={handleViewEvent}
            />
          )}

          {totalPages > 1 && (
            <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
              <Pagination
                count={totalPages}
                page={page}
                onChange={(_, newPage) => setPage(newPage)}
                color="primary"
              />
            </Box>
          )}
        </Paper>

        <CreateEventForm
          open={createFormOpen}
          onClose={() => setCreateFormOpen(false)}
          onSubmit={handleCreateEvent}
        />

        <EventDetailDialog
          open={detailDialogOpen}
          event={selectedEvent}
          onClose={() => {
            setDetailDialogOpen(false);
            setSelectedEvent(null);
          }}
        />

        <Snackbar
          open={!!success}
          autoHideDuration={6000}
          onClose={() => setSuccess(null)}
        >
          <Alert severity="success" onClose={() => setSuccess(null)}>
            {success}
          </Alert>
        </Snackbar>
      </Container>
    </>
  );
}
