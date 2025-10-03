"use client";

import React, { useState, useEffect } from "react";
import {
  Container,
  Typography,
  Card,
  CardContent,
  CardHeader,
  Grid,
  Chip,
  Box,
  CircularProgress,
  Alert,
} from "@mui/material";
import { Navigation } from "../../components/Navigation";

interface ServiceHealth {
  name: string;
  url: string;
  status: "healthy" | "unhealthy" | "loading";
  responseTime?: number;
  error?: string;
}

interface SimulationStatus {
  simulation_mode_enabled: boolean;
  force_webhook_failures: boolean;
  disable_publishing: boolean;
  circuit_breaker_demo_mode: boolean;
  partial_failure_mode: boolean;
  simulate_network_delays: boolean;
  circuit_breaker_state: string;
  circuit_failure_count: number;
  circuit_last_failure: string;
}

export default function StatusPage() {
  const [services, setServices] = useState<ServiceHealth[]>([
    {
      name: "Outbox API",
      url: "http://localhost:8080/health",
      status: "loading",
    },
    {
      name: "Feature Flags API",
      url: "http://localhost:4000/health",
      status: "loading",
    },
    {
      name: "Observability",
      url: "http://localhost:8081/health",
      status: "loading",
    },
  ]);

  const [simulationStatus, setSimulationStatus] =
    useState<SimulationStatus | null>(null);
  const [systemHealth, setSystemHealth] = useState<
    "healthy" | "down" | "loading"
  >("loading");
  const [simulationState, setSimulationState] = useState<
    "normal" | "ready" | "simulating"
  >("normal");

  const checkServiceHealth = async (
    service: ServiceHealth
  ): Promise<ServiceHealth> => {
    const startTime = Date.now();

    try {
      const response = await fetch(service.url, {
        method: "GET",
        signal: AbortSignal.timeout(5000), // 5 second timeout
      });

      const responseTime = Date.now() - startTime;

      if (response.ok) {
        return {
          ...service,
          status: "healthy",
          responseTime,
        };
      } else {
        return {
          ...service,
          status: "unhealthy",
          responseTime,
          error: `HTTP ${response.status}`,
        };
      }
    } catch (error) {
      return {
        ...service,
        status: "unhealthy",
        responseTime: Date.now() - startTime,
        error: error instanceof Error ? error.message : "Unknown error",
      };
    }
  };

  const fetchSimulationStatus = async () => {
    try {
      const response = await fetch("/api/outbox/admin/simulation-status");
      if (response.ok) {
        const data = await response.json();
        setSimulationStatus(data.simulation_status);
      }
    } catch (error) {
      console.error("Failed to fetch simulation status:", error);
    }
  };

  const calculateSystemHealth = (services: ServiceHealth[]) => {
    const unhealthyServices = services.filter((s) => s.status === "unhealthy");

    if (unhealthyServices.length > 0) {
      return "down";
    } else {
      return "healthy";
    }
  };

  const getSimulationState = (simStatus: SimulationStatus | null) => {
    if (!simStatus?.simulation_mode_enabled) {
      return "normal";
    }

    const activeSimulations =
      Object.values(simStatus).filter(Boolean).length - 1; // -1 for simulation_mode_enabled
    return activeSimulations > 0 ? "simulating" : "ready";
  };

  useEffect(() => {
    const checkAllServices = async () => {
      const healthChecks = services.map(checkServiceHealth);
      const results = await Promise.all(healthChecks);
      setServices(results);
    };

    checkAllServices();
    fetchSimulationStatus();

    // Refresh every 30 seconds
    const interval = setInterval(() => {
      checkAllServices();
      fetchSimulationStatus();
    }, 30000);

    return () => clearInterval(interval);
  }, [services]);

  useEffect(() => {
    setSystemHealth(calculateSystemHealth(services));
    setSimulationState(getSimulationState(simulationStatus));
  }, [services, simulationStatus]);

  const getHealthColor = (status: string) => {
    switch (status) {
      case "healthy":
        return "success";
      case "down":
        return "error";
      default:
        return "default";
    }
  };

  const getSimulationColor = (state: string) => {
    switch (state) {
      case "normal":
        return "default";
      case "ready":
        return "info";
      case "simulating":
        return "warning";
      default:
        return "default";
    }
  };

  return (
    <>
      <Navigation />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          System Status
        </Typography>

        {/* System Health & Simulation Status */}
        <Grid container spacing={3} sx={{ mb: 4 }}>
          <Grid size={{ xs: 12, md: 6 }}>
            <Card>
              <CardHeader
                title={
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <Typography variant="h6">System Health</Typography>
                    <Chip
                      label={systemHealth.toUpperCase()}
                      color={getHealthColor(systemHealth)}
                      size="medium"
                    />
                  </Box>
                }
              />
            </Card>
          </Grid>
          <Grid size={{ xs: 12, md: 6 }}>
            <Card>
              <CardHeader
                title={
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <Typography variant="h6">Simulation Mode</Typography>
                    <Chip
                      label={simulationState.toUpperCase()}
                      color={getSimulationColor(simulationState)}
                      size="medium"
                    />
                  </Box>
                }
              />
            </Card>
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          {/* Service Health */}
          <Grid size={{ xs: 12, md: 6 }}>
            <Card>
              <CardHeader title="Service Health" />
              <CardContent>
                {services.map((service) => (
                  <Box
                    key={service.name}
                    sx={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center",
                      py: 1,
                      px: 1,
                      borderBottom: "1px solid #eee",
                      minWidth: 300,
                    }}
                  >
                    <Typography variant="body1">{service.name}</Typography>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      {service.status === "loading" ? (
                        <CircularProgress size={16} />
                      ) : (
                        <>
                          <Chip
                            label={service.status}
                            color={
                              service.status === "healthy" ? "success" : "error"
                            }
                            size="small"
                          />
                          {service.responseTime && (
                            <Typography variant="caption" color="textSecondary">
                              {service.responseTime}ms
                            </Typography>
                          )}
                        </>
                      )}
                    </Box>
                  </Box>
                ))}
              </CardContent>
            </Card>
          </Grid>

          {/* Simulation Status */}
          <Grid size={{ xs: 12, md: 6 }}>
            <Card>
              <CardHeader title="Simulation Status" />
              <CardContent>
                {simulationStatus ? (
                  <Box
                    sx={{ display: "flex", flexDirection: "column", gap: 1 }}
                  >
                    <Box
                      sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                      }}
                    >
                      <Typography variant="body1">Simulation Mode</Typography>
                      <Chip
                        label={
                          simulationStatus.simulation_mode_enabled
                            ? "ENABLED"
                            : "DISABLED"
                        }
                        color={
                          simulationStatus.simulation_mode_enabled
                            ? "warning"
                            : "default"
                        }
                        size="small"
                      />
                    </Box>

                    {simulationStatus.simulation_mode_enabled && (
                      <>
                        {simulationStatus.circuit_breaker_demo_mode && (
                          <Box
                            sx={{
                              mt: 2,
                              p: 2,
                              border: "1px solid #ddd",
                              borderRadius: 1,
                            }}
                          >
                            <Typography variant="subtitle2" gutterBottom>
                              Circuit Breaker Status
                            </Typography>
                            <Box
                              sx={{
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                                mb: 1,
                              }}
                            >
                              <Typography variant="body2">State:</Typography>
                              <Chip
                                label={simulationStatus.circuit_breaker_state}
                                color={
                                  simulationStatus.circuit_breaker_state ===
                                  "CLOSED"
                                    ? "success"
                                    : simulationStatus.circuit_breaker_state ===
                                      "OPEN"
                                    ? "error"
                                    : "warning"
                                }
                                size="small"
                              />
                            </Box>
                            <Box
                              sx={{
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                              }}
                            >
                              <Typography variant="body2">Failures:</Typography>
                              <Typography variant="body2">
                                {simulationStatus.circuit_failure_count}
                              </Typography>
                            </Box>
                          </Box>
                        )}

                        {simulationStatus.disable_publishing && (
                          <Alert severity="warning" sx={{ mt: 1 }}>
                            Publishing is disabled - events will remain pending
                          </Alert>
                        )}
                        {simulationStatus.force_webhook_failures && (
                          <Alert severity="error" sx={{ mt: 1 }}>
                            Webhook failures are being forced
                          </Alert>
                        )}
                        {simulationStatus.partial_failure_mode && (
                          <Alert severity="info" sx={{ mt: 1 }}>
                            Partial failure mode is active
                          </Alert>
                        )}
                        {simulationStatus.simulate_network_delays && (
                          <Alert severity="info" sx={{ mt: 1 }}>
                            Network delays are being simulated
                          </Alert>
                        )}
                      </>
                    )}
                  </Box>
                ) : (
                  <Typography color="textSecondary">
                    Loading simulation status...
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Container>
    </>
  );
}
