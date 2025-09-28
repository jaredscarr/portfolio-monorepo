import { Typography, Box, Container, Button, Chip } from "@mui/material";
import { Navigation } from "../../../components/Navigation";
import { ArrowBack } from "@mui/icons-material";
import Link from "next/link";
import { notFound } from "next/navigation";

const caseStudyData: Record<
  string,
  {
    title: string;
    description: string;
    technologies: string[];
    status: string;
    overview: string;
    keyFeatures: string[];
    challenges: string[];
    learnings: string[];
    repositoryUrl: string;
    docsUrl: string;
    liveApiUrl: string;
  }
> = {
  "outbox-pattern": {
    title: "Outbox Pattern Implementation",
    description:
      "A practical implementation of the transactional outbox pattern for reliable event publishing in distributed systems.",
    technologies: ["Go", "PostgreSQL", "Event Sourcing", "Docker"],
    status: "Completed",
    overview: `
      The transactional outbox pattern is a reliable way to publish events in distributed systems 
      while maintaining data consistency. This implementation demonstrates how to ensure that 
      database changes and event publishing happen atomically, preventing data inconsistencies 
      that can occur when external systems are unavailable.
    `,
    keyFeatures: [
      "Atomic database transactions with event storage",
      "Reliable event publishing with retry mechanisms",
      "Event ordering and deduplication",
      "Health monitoring and observability",
      "Docker containerization for easy deployment",
    ],
    challenges: [
      "Ensuring exactly-once event delivery",
      "Managing event ordering across multiple aggregates",
      "Handling system failures and recovery",
      "Monitoring and alerting for event processing delays",
    ],
    learnings: [
      "Deep understanding of distributed system challenges",
      "Event-driven architecture patterns",
      "Database transaction management",
      "System observability and monitoring",
    ],
    repositoryUrl: "https://github.com/yourusername/outbox-api",
    docsUrl: "http://localhost:8080/swagger/index.html",
    liveApiUrl: "/outbox",
  },
  "feature-flags-api": {
    title: "Feature Flags Service",
    description:
      "A lightweight feature flags service with dynamic configuration and A/B testing capabilities.",
    technologies: ["Go", "JSON Config", "REST API", "Swagger"],
    status: "In Progress",
    overview: `
      A flexible feature flags service that allows for dynamic feature toggling without 
      code deployments. Supports percentage-based rollouts, user targeting, and A/B testing 
      scenarios for controlled feature releases.
    `,
    keyFeatures: [
      "Dynamic feature flag configuration",
      "Percentage-based feature rollouts",
      "User and group targeting",
      "RESTful API with comprehensive documentation",
      "Environment-specific configurations",
    ],
    challenges: [
      "Designing a flexible configuration schema",
      "Implementing efficient flag evaluation",
      "Ensuring high availability and low latency",
      "Managing configuration changes safely",
    ],
    learnings: [
      "API design principles",
      "Configuration management strategies",
      "Performance optimization techniques",
      "Documentation and API usability",
    ],
    repositoryUrl: "https://github.com/yourusername/feature-flags-api",
    docsUrl: "http://localhost:4000/swagger/index.html",
    liveApiUrl: "http://localhost:4000",
  },
  observability: {
    title: "Observability",
    description:
      "A standardized observability service providing health checks, readiness probes, and Prometheus metrics exposure.",
    technologies: ["Go", "Prometheus", "Swagger", "OpenAPI", "Unit Testing"],
    status: "Completed",
    overview: `
      The Metrics API provides basic observability endpoints for services in this portfolio, 
      standardizing health checks, readiness checks, and Prometheus metrics exposure. This 
      service establishes a minimal baseline for monitoring with health, readiness, HTTP 
      counters, latency histograms, and Go runtime metrics.
    `,
    keyFeatures: [
      "Standardized health and readiness probe endpoints",
      "Prometheus metrics exposure in text format",
      "HTTP request counters and latency histograms",
      "Go runtime metrics monitoring",
      "Comprehensive Swagger/OpenAPI documentation",
      "Extensive unit test coverage with benchmarks",
    ],
    challenges: [
      "Designing a flexible metrics registration system",
      "Ensuring consistent observability across all services",
      "Balancing minimal overhead with comprehensive monitoring",
      "Creating reusable patterns for future service integration",
    ],
    learnings: [
      "Observability best practices in distributed systems",
      "Prometheus metrics design and implementation",
      "Comprehensive testing strategies including benchmarks",
      "API documentation with Swagger/OpenAPI standards",
      "Go runtime monitoring and performance analysis",
    ],
    repositoryUrl: "https://github.com/yourusername/observability",
    docsUrl: "http://localhost:8081/docs/index.html",
    liveApiUrl: "http://localhost:8081",
  },
};

interface PageProps {
  params: Promise<{
    slug: string;
  }>;
}

async function getCaseStudy(params: Promise<{ slug: string }>) {
  const { slug } = await params;
  const study = caseStudyData[slug];

  if (!study) {
    notFound();
  }

  return study;
}

export default async function CaseStudyDetail({ params }: PageProps) {
  const study = await getCaseStudy(params);

  return (
    <>
      <Navigation />

      <Container maxWidth="lg">
        <Box sx={{ py: 4 }}>
          <Button
            component={Link}
            href="/case-studies"
            startIcon={<ArrowBack />}
            sx={{ mb: 3 }}
          >
            Back to Case Studies
          </Button>

          <Typography variant="h3" component="h1" gutterBottom>
            {study.title}
          </Typography>

          <Box sx={{ display: "flex", gap: 1, mb: 2, flexWrap: "wrap" }}>
            <Chip
              label={`Status: ${study.status}`}
              color="primary"
              variant="outlined"
            />
            {study.technologies.map((tech: string) => (
              <Chip key={tech} label={tech} variant="outlined" />
            ))}
          </Box>

          <Typography variant="h5" color="text.secondary" sx={{ mb: 4 }}>
            {study.description}
          </Typography>

          <Typography variant="h4" component="h2" gutterBottom sx={{ mt: 4 }}>
            Overview
          </Typography>
          <Typography variant="body1" sx={{ mb: 4 }}>
            {study.overview}
          </Typography>

          <Typography variant="h4" component="h2" gutterBottom>
            Key Features
          </Typography>
          <Box component="ul" sx={{ mb: 4 }}>
            {study.keyFeatures.map((feature: string, index: number) => (
              <Typography component="li" key={index} sx={{ mb: 1 }}>
                {feature}
              </Typography>
            ))}
          </Box>

          <Typography variant="h4" component="h2" gutterBottom>
            Technical Challenges
          </Typography>
          <Box component="ul" sx={{ mb: 4 }}>
            {study.challenges.map((challenge: string, index: number) => (
              <Typography component="li" key={index} sx={{ mb: 1 }}>
                {challenge}
              </Typography>
            ))}
          </Box>

          <Typography variant="h4" component="h2" gutterBottom>
            Key Learnings
          </Typography>
          <Box component="ul" sx={{ mb: 4 }}>
            {study.learnings.map((learning: string, index: number) => (
              <Typography component="li" key={index} sx={{ mb: 1 }}>
                {learning}
              </Typography>
            ))}
          </Box>

          <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap", mt: 4 }}>
            <Button
              variant="contained"
              href={study.repositoryUrl}
              target="_blank"
            >
              View Repository
            </Button>
            <Button variant="outlined" href={study.docsUrl} target="_blank">
              Swagger Docs
            </Button>
            <Button variant="outlined" href={study.liveApiUrl} target="_blank">
              Live API
            </Button>
          </Box>
        </Box>
      </Container>
    </>
  );
}
