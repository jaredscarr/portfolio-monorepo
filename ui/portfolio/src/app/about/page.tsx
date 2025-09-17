"use client";

import { Typography, Box, Container } from "@mui/material";
import { Navigation } from "../../components/Navigation";

export default function About() {
  return (
    <>
      <Navigation />

      <Container maxWidth="lg">
        <Box sx={{ py: 4 }}>
          <Typography variant="h3" component="h1" gutterBottom>
            About Me
          </Typography>
          <Typography variant="h5" color="text.secondary" sx={{ mb: 4 }}>
            Software Engineer & System Design Enthusiast
          </Typography>
          <Typography variant="body1" sx={{ mb: 3 }}>
            I'm passionate about exploring modern software architecture patterns
            and distributed systems design. This portfolio represents my journey
            of learning and experimenting with various technologies and
            architectural approaches.
          </Typography>
          <Typography variant="body1" sx={{ mb: 3 }}>
            Based in the Pacific Northwest, I enjoy diving deep into technical
            challenges and sharing my findings through practical implementations
            and case studies.
          </Typography>
          <Typography variant="body1">
            Each case study in this portfolio represents a hands-on exploration
            of a specific pattern, technology, or architectural concept that I
            found interesting and worth investigating further.
          </Typography>
        </Box>
      </Container>
    </>
  );
}
