"use client";

import React, { useState, useEffect } from "react";
import {
  Card,
  CardContent,
  CardHeader,
  FormControlLabel,
  Switch,
  Typography,
  Box,
  Chip,
  Alert,
} from "@mui/material";

interface SimulationStatus {
  simulation_mode_enabled: boolean;
  force_webhook_failures: boolean;
  disable_publishing: boolean;
  circuit_breaker_demo_mode: boolean;
  partial_failure_mode: boolean;
  simulate_network_delays: boolean;
}

interface SimulationControlsProps {
  onSimulationChange?: () => void;
}

export const SimulationControls: React.FC<SimulationControlsProps> = ({
  onSimulationChange,
}) => {
  const [simulationStatus, setSimulationStatus] =
    useState<SimulationStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch current simulation status
  const fetchSimulationStatus = async () => {
    try {
      const response = await fetch("/api/outbox/admin/simulation-status");
      if (!response.ok) {
        throw new Error("Failed to fetch simulation status");
      }
      const data = await response.json();
      setSimulationStatus(data.simulation_status);
    } catch (err) {
      console.error("Failed to fetch simulation status:", err);
      setError("Failed to load simulation status");
    }
  };

  // Update a specific flag
  const updateFlag = async (flagKey: string, enabled: boolean) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/feature-flags/admin/flags/${flagKey}?env=local`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ enabled }),
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to update ${flagKey}`);
      }

      // Refresh simulation status
      await fetchSimulationStatus();

      // Notify parent component
      onSimulationChange?.();
    } catch (err) {
      console.error(`Failed to update ${flagKey}:`, err);
      setError(`Failed to update ${flagKey}`);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSimulationStatus();
  }, []);

  if (!simulationStatus) {
    return (
      <Card>
        <CardHeader title="Simulation Controls" />
        <CardContent>
          <Typography>Loading simulation status...</Typography>
        </CardContent>
      </Card>
    );
  }

  const simulationFlags = [
    {
      key: "disable_publishing",
      label: "Disable Publishing",
      description: "Keep events in pending state",
      color: "warning" as const,
    },
    {
      key: "force_webhook_failures",
      label: "Force Webhook Failures",
      description: "All events fail with simulated error",
      color: "error" as const,
    },
    {
      key: "partial_failure_mode",
      label: "Partial Failures",
      description: "Mix of success and failure (every 3rd fails)",
      color: "info" as const,
    },
    {
      key: "simulate_network_delays",
      label: "Network Delays",
      description: "Add 2-second delays to simulate slow responses",
      color: "secondary" as const,
    },
    {
      key: "circuit_breaker_demo_mode",
      label: "Circuit Breaker Demo",
      description: "Demonstrate circuit breaker behavior",
      color: "primary" as const,
    },
  ];

  return (
    <Card>
      <CardHeader
        title={
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <Typography variant="h6">🎛️ Simulation Controls</Typography>
            <Chip
              label={
                simulationStatus.simulation_mode_enabled
                  ? "ENABLED"
                  : "DISABLED"
              }
              color={
                simulationStatus.simulation_mode_enabled ? "success" : "default"
              }
              size="small"
            />
          </Box>
        }
      />
      <CardContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {!simulationStatus.simulation_mode_enabled && (
          <Alert severity="info" sx={{ mb: 2 }}>
            Simulation mode is disabled. Enable it to use simulation controls.
          </Alert>
        )}

        <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
          {/* Master simulation mode toggle */}
          <FormControlLabel
            control={
              <Switch
                checked={simulationStatus.simulation_mode_enabled}
                onChange={(e) =>
                  updateFlag("simulation_mode_enabled", e.target.checked)
                }
                disabled={loading}
              />
            }
            label={
              <Box>
                <Typography variant="body1" fontWeight="bold">
                  Master Simulation Mode
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Enable/disable all simulation features
                </Typography>
              </Box>
            }
          />

          {/* Individual simulation controls */}
          {simulationFlags.map((flag) => (
            <FormControlLabel
              key={flag.key}
              control={
                <Switch
                  checked={
                    simulationStatus[
                      flag.key as keyof SimulationStatus
                    ] as boolean
                  }
                  onChange={(e) => updateFlag(flag.key, e.target.checked)}
                  disabled={
                    loading || !simulationStatus.simulation_mode_enabled
                  }
                />
              }
              label={
                <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                  <Box>
                    <Typography variant="body1" fontWeight="medium">
                      {flag.label}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {flag.description}
                    </Typography>
                  </Box>
                  {simulationStatus[flag.key as keyof SimulationStatus] && (
                    <Chip label="ACTIVE" color={flag.color} size="small" />
                  )}
                </Box>
              }
            />
          ))}
        </Box>
      </CardContent>
    </Card>
  );
};
