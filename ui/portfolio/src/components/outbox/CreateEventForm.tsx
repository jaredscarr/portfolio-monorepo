"use client";

import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Typography,
  Alert,
} from "@mui/material";
import { CreateEventRequest } from "../../types/outbox";

interface CreateEventFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (event: CreateEventRequest) => Promise<void>;
}

export const CreateEventForm: React.FC<CreateEventFormProps> = ({
  open,
  onClose,
  onSubmit,
}) => {
  const [formData, setFormData] = useState<CreateEventRequest>({
    type: "",
    source: "",
    data: {},
    metadata: {},
  });
  const [dataJson, setDataJson] = useState("{}");
  const [metadataJson, setMetadataJson] = useState("{}");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async () => {
    try {
      setError(null);
      setLoading(true);

      // Validate JSON
      const data = JSON.parse(dataJson);
      const metadata = JSON.parse(metadataJson);

      await onSubmit({
        ...formData,
        data,
        metadata,
      });

      // Reset form
      setFormData({
        type: "",
        source: "",
        data: {},
        metadata: {},
      });
      setDataJson("{}");
      setMetadataJson("{}");
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create event");
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>Create New Event</DialogTitle>
      <DialogContent>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2, mt: 1 }}>
          {error && (
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          <TextField
            label="Event Type"
            value={formData.type}
            onChange={(e) => setFormData({ ...formData, type: e.target.value })}
            fullWidth
            required
            placeholder="e.g., user.created, order.completed"
          />

          <TextField
            label="Source"
            value={formData.source}
            onChange={(e) =>
              setFormData({ ...formData, source: e.target.value })
            }
            fullWidth
            required
            placeholder="e.g., user-service, order-service"
          />

          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Event Data (JSON)
            </Typography>
            <TextField
              value={dataJson}
              onChange={(e) => setDataJson(e.target.value)}
              multiline
              rows={4}
              fullWidth
              placeholder='{"user_id": "123", "email": "user@example.com"}'
              error={(() => {
                try {
                  JSON.parse(dataJson);
                  return false;
                } catch {
                  return true;
                }
              })()}
              helperText={(() => {
                try {
                  JSON.parse(dataJson);
                  return "Valid JSON";
                } catch {
                  return "Invalid JSON format";
                }
              })()}
            />
          </Box>

          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Metadata (JSON) - Optional
            </Typography>
            <TextField
              value={metadataJson}
              onChange={(e) => setMetadataJson(e.target.value)}
              multiline
              rows={3}
              fullWidth
              placeholder='{"version": "1.0", "correlation_id": "abc-123"}'
              error={(() => {
                try {
                  JSON.parse(metadataJson);
                  return false;
                } catch {
                  return true;
                }
              })()}
              helperText={(() => {
                try {
                  JSON.parse(metadataJson);
                  return "Valid JSON";
                } catch {
                  return "Invalid JSON format";
                }
              })()}
            />
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={
            loading ||
            !formData.type ||
            !formData.source ||
            (() => {
              try {
                JSON.parse(dataJson);
                JSON.parse(metadataJson);
                return false;
              } catch {
                return true;
              }
            })()
          }
        >
          {loading ? "Creating..." : "Create Event"}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
