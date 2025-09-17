"use client";

import { Typography, Button, Box, Container } from "@mui/material";
import { Navigation } from "../components/Navigation";

export default function Home() {
  return (
    <>
      <Navigation />

      <Container maxWidth="lg">
        <Box sx={{ py: 4 }}>
          <Typography variant="h3" component="h1" gutterBottom>
            Welcome to My Portfolio
          </Typography>
          <Typography variant="h5" color="text.secondary" sx={{ mb: 4 }}>
            Exploring software architecture and system design through hands-on
            case studies
          </Typography>
          <Typography variant="body1" sx={{ mb: 3 }}>
            This portfolio showcases my technical explorations and learning
            journey through various software engineering patterns and practices.
          </Typography>
          <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
            <Button variant="contained" color="primary" href="/case-studies">
              View Case Studies
            </Button>
            <Button variant="outlined" color="secondary" href="/about">
              About Me
            </Button>
          </Box>
        </Box>
      </Container>
    </>
  );
}
