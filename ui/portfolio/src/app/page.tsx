"use client";

import { Typography, Button } from "@mui/material";

export default function Home() {
  return (
    <div>
      <Typography variant="h4">Welcome to Your Portfolio</Typography>
      <Typography variant="body1">
        This is a minimal test without sx props.
      </Typography>
      <Button variant="contained" color="primary">
        Test Button
      </Button>
    </div>
  );
}
