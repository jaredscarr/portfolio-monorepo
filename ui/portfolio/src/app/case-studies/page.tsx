"use client";

import {
  Typography,
  Box,
  Container,
  Card,
  CardContent,
  Button,
  Grid,
} from "@mui/material";
import { Navigation } from "../../components/Navigation";
import Link from "next/link";

const caseStudies = [
  {
    slug: "outbox-pattern",
    title: "Outbox Pattern Implementation",
    description:
      "A practical implementation of the transactional outbox pattern for reliable event publishing in distributed systems.",
    technologies: ["Go", "PostgreSQL", "Event Sourcing"],
    status: "Completed",
  },
  {
    slug: "feature-flags-api",
    title: "Feature Flags Service",
    description:
      "A lightweight feature flags service with dynamic configuration and A/B testing capabilities.",
    technologies: ["Go", "JSON Config", "REST API"],
    status: "In Progress",
  },
];

export default function CaseStudies() {
  return (
    <>
      <Navigation />

      <Container maxWidth="lg">
        <Box sx={{ py: 4 }}>
          <Typography variant="h3" component="h1" gutterBottom>
            Case Studies
          </Typography>
          <Typography variant="h5" color="text.secondary" sx={{ mb: 4 }}>
            Technical explorations and architectural experiments
          </Typography>
          <Typography variant="body1" sx={{ mb: 4 }}>
            Each case study represents a deep dive into a specific software
            engineering pattern or technology. These are learning projects that
            explore real-world challenges and solutions.
          </Typography>

          <Grid container spacing={3}>
            {caseStudies.map((study) => (
              <Grid size={{ xs: 12, md: 6 }} key={study.slug}>
                <Card
                  sx={{
                    height: "100%",
                    display: "flex",
                    flexDirection: "column",
                  }}
                >
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Typography variant="h5" component="h2" gutterBottom>
                      {study.title}
                    </Typography>
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ mb: 2 }}
                    >
                      Status: {study.status}
                    </Typography>
                    <Typography variant="body1" sx={{ mb: 2 }}>
                      {study.description}
                    </Typography>
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ mb: 3 }}
                    >
                      Technologies: {study.technologies.join(", ")}
                    </Typography>
                    <Button
                      variant="contained"
                      component={Link}
                      href={`/case-studies/${study.slug}`}
                      sx={{ mt: "auto" }}
                    >
                      View Case Study
                    </Button>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Box>
      </Container>
    </>
  );
}
