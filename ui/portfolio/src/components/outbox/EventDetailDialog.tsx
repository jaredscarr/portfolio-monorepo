"use client";

import React from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Chip,
  Divider,
  Paper,
} from "@mui/material";
import { OutboxEvent } from "../../types/outbox";

interface EventDetailDialogProps {
  open: boolean;
  event: OutboxEvent | null;
  onClose: () => void;
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

export const EventDetailDialog: React.FC<EventDetailDialogProps> = ({
  open,
  event,
  onClose,
}) => {
  if (!event) return null;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Event Details</DialogTitle>
      <DialogContent>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
          <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
            <Typography variant="h6">Status:</Typography>
            <Chip
              label={event.status}
              color={getStatusColor(event.status)}
              size="medium"
            />
          </Box>

          <Box>
            <Typography variant="subtitle1" fontWeight="bold">
              Basic Information
            </Typography>
            <Paper sx={{ p: 2, mt: 1, bgcolor: "grey.50" }}>
              <Typography variant="body2">
                <strong>ID:</strong> {event.id}
              </Typography>
              <Typography variant="body2">
                <strong>Type:</strong> {event.type}
              </Typography>
              <Typography variant="body2">
                <strong>Source:</strong> {event.source}
              </Typography>
              <Typography variant="body2">
                <strong>Retry Count:</strong> {event.retry_count}
              </Typography>
              <Typography variant="body2">
                <strong>Created At:</strong>{" "}
                {new Date(event.created_at).toLocaleString()}
              </Typography>
              {event.published_at && (
                <Typography variant="body2">
                  <strong>Published At:</strong>{" "}
                  {new Date(event.published_at).toLocaleString()}
                </Typography>
              )}
            </Paper>
          </Box>

          {event.error_message && (
            <Box>
              <Typography variant="subtitle1" fontWeight="bold" color="error">
                Error Message
              </Typography>
              <Paper
                sx={{
                  p: 2,
                  mt: 1,
                  bgcolor: "error.light",
                  color: "error.contrastText",
                }}
              >
                <Typography variant="body2" fontFamily="monospace">
                  {event.error_message}
                </Typography>
              </Paper>
            </Box>
          )}

          <Divider />

          <Box>
            <Typography variant="subtitle1" fontWeight="bold">
              Event Data
            </Typography>
            <Paper sx={{ p: 2, mt: 1, bgcolor: "grey.50" }}>
              <pre
                style={{
                  margin: 0,
                  fontFamily: "monospace",
                  fontSize: "12px",
                  whiteSpace: "pre-wrap",
                  wordBreak: "break-word",
                }}
              >
                {JSON.stringify(event.data, null, 2)}
              </pre>
            </Paper>
          </Box>

          {event.metadata && Object.keys(event.metadata).length > 0 && (
            <Box>
              <Typography variant="subtitle1" fontWeight="bold">
                Metadata
              </Typography>
              <Paper sx={{ p: 2, mt: 1, bgcolor: "grey.50" }}>
                <pre
                  style={{
                    margin: 0,
                    fontFamily: "monospace",
                    fontSize: "12px",
                    whiteSpace: "pre-wrap",
                    wordBreak: "break-word",
                  }}
                >
                  {JSON.stringify(event.metadata, null, 2)}
                </pre>
              </Paper>
            </Box>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};
