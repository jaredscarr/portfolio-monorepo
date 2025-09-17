"use client";

import { Typography, Button, Box, AppBar, Toolbar } from "@mui/material";
import { ThemeToggle } from "../components/ThemeToggle";

export default function Home() {
  return (
    <>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Portfolio
          </Typography>
          <ThemeToggle />
        </Toolbar>
      </AppBar>

      <Box sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          Welcome to Your Portfolio
        </Typography>
        <Typography variant="body1" sx={{ mb: 2 }}>
          This portfolio now supports light and dark themes! Use the toggle in
          the top-right corner.
        </Typography>
        <Button variant="contained" color="primary" sx={{ mr: 2 }}>
          Primary Button
        </Button>
        <Button variant="outlined" color="secondary">
          Secondary Button
        </Button>
      </Box>
    </>
  );
}
