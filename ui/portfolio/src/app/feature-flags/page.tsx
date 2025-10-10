"use client";

import React, { useState, useEffect } from "react";
import {
  Container,
  Typography,
  Box,
  Paper,
  Alert,
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  ToggleButtonGroup,
  ToggleButton,
  Switch,
  Snackbar,
  Button,
  Card,
  CardContent,
  TextField,
  InputAdornment,
} from "@mui/material";
import {
  Refresh as RefreshIcon,
  RestartAlt as RestartAltIcon,
  Search as SearchIcon,
} from "@mui/icons-material";
import { Navigation } from "../../components/Navigation";
import { FlagsResponse, UpdateFlagRequest } from "../../types/feature-flags";

export default function FeatureFlagsPage() {
  const [flags, setFlags] = useState<FlagsResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [environment, setEnvironment] = useState<"local" | "prod">("local");
  const [updatingFlags, setUpdatingFlags] = useState<Set<string>>(new Set());
  const [searchQuery, setSearchQuery] = useState("");

  const fetchFlags = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(
        `/api/feature-flags/flags?env=${environment}`
      );

      if (!response.ok) {
        throw new Error("Failed to fetch flags");
      }

      const data: FlagsResponse = await response.json();
      setFlags(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch flags");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFlags();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [environment]);

  const handleEnvironmentChange = (
    _event: React.MouseEvent<HTMLElement>,
    newEnvironment: "local" | "prod" | null
  ) => {
    if (newEnvironment !== null) {
      setEnvironment(newEnvironment);
    }
  };

  const handleToggleFlag = async (flagKey: string, currentValue: boolean) => {
    try {
      // Add flag to updating set
      setUpdatingFlags((prev) => new Set(prev).add(flagKey));
      setError(null);

      const updateRequest: UpdateFlagRequest = {
        enabled: !currentValue,
      };

      const response = await fetch(
        `/api/feature-flags/admin/flags/${flagKey}?env=${environment}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(updateRequest),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to update flag");
      }

      const result = await response.json();

      // Update local state
      setFlags((prevFlags) => {
        if (!prevFlags) return prevFlags;
        return {
          ...prevFlags,
          [flagKey]: result.enabled,
        };
      });

      setSuccess(
        `Flag "${flagKey}" ${
          result.enabled ? "enabled" : "disabled"
        } successfully`
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update flag");
    } finally {
      // Remove flag from updating set
      setUpdatingFlags((prev) => {
        const newSet = new Set(prev);
        newSet.delete(flagKey);
        return newSet;
      });
    }
  };

  const handleReloadFlags = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch("/api/feature-flags/admin/reload", {
        method: "POST",
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Failed to reload flags");
      }

      // Refetch flags after reload
      await fetchFlags();
      setSuccess("Flags reloaded from disk successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to reload flags");
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = () => {
    fetchFlags();
  };

  const allFlagsArray = flags
    ? Object.entries(flags).map(([key, enabled]) => ({ key, enabled }))
    : [];

  const flagsArray = allFlagsArray.filter((flag) =>
    flag.key.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const enabledCount = flagsArray.filter((flag) => flag.enabled).length;
  const disabledCount = flagsArray.filter((flag) => !flag.enabled).length;

  return (
    <>
      <Navigation />
      <Container maxWidth="xl" sx={{ py: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom>
            Feature Flags
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Manage and monitor feature flags across different environments
          </Typography>
        </Box>

        {error && (
          <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {/* Stats Cards */}
        <Box
          sx={{
            display: "grid",
            gridTemplateColumns: {
              xs: "1fr",
              sm: "repeat(3, 1fr)",
            },
            gap: 3,
            mb: 3,
          }}
        >
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom variant="body2">
                Total Flags
              </Typography>
              <Typography variant="h4">{allFlagsArray.length}</Typography>
            </CardContent>
          </Card>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom variant="body2">
                Enabled
              </Typography>
              <Typography variant="h4" color="success.main">
                {enabledCount}
              </Typography>
            </CardContent>
          </Card>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom variant="body2">
                Disabled
              </Typography>
              <Typography variant="h4" color="text.secondary">
                {disabledCount}
              </Typography>
            </CardContent>
          </Card>
        </Box>

        <Paper sx={{ p: 3 }}>
          <Box
            sx={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              mb: 3,
              flexWrap: "wrap",
              gap: 2,
            }}
          >
            <Box
              sx={{
                display: "flex",
                gap: 2,
                alignItems: "center",
                flexWrap: "wrap",
              }}
            >
              <Typography variant="h6">Flags</Typography>
              <Button
                variant="outlined"
                size="small"
                startIcon={<RefreshIcon />}
                onClick={handleRefresh}
                disabled={loading}
              >
                Refresh
              </Button>
              <Button
                variant="outlined"
                size="small"
                startIcon={<RestartAltIcon />}
                onClick={handleReloadFlags}
                disabled={loading}
                color="warning"
              >
                Reload from Disk
              </Button>
            </Box>

            <ToggleButtonGroup
              value={environment}
              exclusive
              onChange={handleEnvironmentChange}
              aria-label="environment"
              size="small"
            >
              <ToggleButton value="local" aria-label="local environment">
                Local
              </ToggleButton>
              <ToggleButton value="prod" aria-label="production environment">
                Production
              </ToggleButton>
            </ToggleButtonGroup>
          </Box>

          <Box sx={{ mb: 3 }}>
            <TextField
              fullWidth
              size="small"
              placeholder="Search flags by name..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
            {searchQuery && (
              <Typography
                variant="caption"
                color="text.secondary"
                sx={{ mt: 1, display: "block" }}
              >
                Showing {flagsArray.length} of {allFlagsArray.length} flags
              </Typography>
            )}
          </Box>

          {loading ? (
            <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
              <CircularProgress />
            </Box>
          ) : flags && flagsArray.length > 0 ? (
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Flag Key</TableCell>
                    <TableCell align="center">Status</TableCell>
                    <TableCell align="center">Toggle</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {flagsArray.map((flag) => {
                    const isUpdating = updatingFlags.has(flag.key);
                    return (
                      <TableRow
                        key={flag.key}
                        sx={{
                          "&:last-child td, &:last-child th": { border: 0 },
                        }}
                      >
                        <TableCell component="th" scope="row">
                          <Typography
                            variant="body2"
                            sx={{ fontFamily: "monospace" }}
                          >
                            {flag.key}
                          </Typography>
                        </TableCell>
                        <TableCell align="center">
                          <Chip
                            label={flag.enabled ? "Enabled" : "Disabled"}
                            color={flag.enabled ? "success" : "default"}
                            size="small"
                          />
                        </TableCell>
                        <TableCell align="center">
                          <Switch
                            checked={flag.enabled}
                            onChange={() =>
                              handleToggleFlag(flag.key, flag.enabled)
                            }
                            disabled={isUpdating}
                            color="primary"
                          />
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </TableContainer>
          ) : (
            <Box sx={{ textAlign: "center", py: 4 }}>
              <Typography color="text.secondary">
                {searchQuery
                  ? `No flags found matching "${searchQuery}"`
                  : "No flags available"}
              </Typography>
              {searchQuery && (
                <Button
                  size="small"
                  onClick={() => setSearchQuery("")}
                  sx={{ mt: 2 }}
                >
                  Clear Search
                </Button>
              )}
            </Box>
          )}
        </Paper>

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
