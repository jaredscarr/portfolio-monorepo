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
            I&apos;m a Full Stack Software Engineer with experience developing
            and maintaining web applications, automation tools, and data
            processing systems.
          </Typography>
          <Typography variant="body1">
            Each case study represents practical implementations of concepts
            I&apos;ve encountered in professional settings, demonstrating my
            commitment to continuous learning and technical excellence.
          </Typography>
        </Box>
      </Container>
    </>
  );
}
