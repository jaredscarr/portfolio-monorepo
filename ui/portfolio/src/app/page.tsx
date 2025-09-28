"use client";

import { Typography, Button, Box, Container } from "@mui/material";
import { Navigation } from "../components/Navigation";

export default function Home() {
  return (
    <>
      <Navigation />
      <Container maxWidth="lg">
        <Box
          sx={{
            py: 4,
            pt: { xs: 4, sm: 8, md: 12, lg: 20 },
            pl: { xs: 2, sm: 4, md: 8, lg: 25 },
          }}
        >
          <Typography
            variant="h1"
            component="h1"
            gutterBottom
            sx={{ fontWeight: "fontWeightBold" }}
          >
            JARED SCARR
          </Typography>
          <Typography variant="h5" color="text.secondary" sx={{ mb: 6 }}>
            Software Engineer, Full Stack, & Aspiring Architect
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
