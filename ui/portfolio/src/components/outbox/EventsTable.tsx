"use client";

import React from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Tooltip,
  Box,
  Typography,
} from "@mui/material";
import {
  Refresh as RefreshIcon,
  Delete as DeleteIcon,
  Visibility as ViewIcon,
} from "@mui/icons-material";
import { OutboxEvent } from "../../types/outbox";

interface EventsTableProps {
  events: OutboxEvent[];
  loading?: boolean;
  onRetry?: (eventId: string) => void;
  onDelete?: (eventId: string) => void;
  onView?: (event: OutboxEvent) => void;
}

const getStatusColor = (status: OutboxEvent["status"]) => {
  switch (status) {
    case "pending":
      return "warning";
    case "published":
      return "success";
    case "failed":
      return "error";
    case "retrying":
      return "info";
    default:
      return "default";
  }
};

export const EventsTable: React.FC<EventsTableProps> = ({
  events,
  loading = false,
  onRetry,
  onDelete,
  onView,
}) => {
  if (events.length === 0) {
    return (
      <Box sx={{ p: 4, textAlign: "center" }}>
        <Typography variant="h6" color="text.secondary">
          No events found
        </Typography>
      </Box>
    );
  }

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Source</TableCell>
            <TableCell>Status</TableCell>
            <TableCell>Retry Count</TableCell>
            <TableCell>Created At</TableCell>
            <TableCell>Actions</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {events.map((event) => (
            <TableRow key={event.id}>
              <TableCell>
                <Typography variant="body2" fontFamily="monospace">
                  {event.id.slice(0, 8)}...
                </Typography>
              </TableCell>
              <TableCell>{event.type}</TableCell>
              <TableCell>{event.source}</TableCell>
              <TableCell>
                <Chip
                  label={event.status}
                  color={getStatusColor(event.status)}
                  size="small"
                />
              </TableCell>
              <TableCell>{event.retry_count}</TableCell>
              <TableCell>
                {new Date(event.created_at).toLocaleDateString()}
              </TableCell>
              <TableCell>
                <Box sx={{ display: "flex", gap: 1 }}>
                  {onView && (
                    <Tooltip title="View details">
                      <IconButton size="small" onClick={() => onView(event)}>
                        <ViewIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                  {onRetry && event.status === "failed" && (
                    <Tooltip title="Retry event">
                      <IconButton
                        size="small"
                        onClick={() => onRetry(event.id)}
                        disabled={loading}
                      >
                        <RefreshIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                  {onDelete && (
                    <Tooltip title="Delete event">
                      <IconButton
                        size="small"
                        onClick={() => onDelete(event.id)}
                        disabled={loading}
                        color="error"
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                </Box>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
