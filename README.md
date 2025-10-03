# Portfolio Monorepo

This repository houses examples of my work as a software engineer. The purpose is twofold:

1. **Learning and practice** — explore Go for API development and deepen my understanding of architecture design patterns and trade-offs.
2. **Showcasing architecture decisions** — demonstrate how to reason about systems, document decisions, and enforce quality through structure and fitness functions.

---

## Why a Monorepo?

* **Shared tooling:** Centralized linting, testing, and observability packages.
* **Consistency:** All services share conventions for health, metrics, logging, and API documentation.
* **Fitness functions:** Easier to enforce cross-cutting architectural rules with tooling (e.g., ESLint for JS, static analysis for Go).
* **Portfolio clarity:** One place to browse multiple services that demonstrate different design patterns.

---

## Microservices Architecture: *Technology Evolution*

While this is a **monorepo** each service is deployed as an **independent microservice**. The decision to use microservices rather than a monolithic runtime architecture was based on:

**Why Microservices:**
* **Technology Evolution:** Each service can adopt different languages, frameworks, or architectural patterns
* **Portfolio Flexibility:** Enables showcasing proficiency across multiple technology stacks (Go, Python, Node.js, etc.)
* **Future-Proofing:** No guarantee all services will remain Go-based as the portfolio evolves
* **Negligible Cost Impact:** Resource difference (~40-100MB) is insignificant on a t3.micro instance

**Trade-offs Accepted:**
* **Operational Complexity:** More containers and service coordination
* **Resource Overhead:** Multiple HTTP servers instead of shared memory
* **Network Latency:** Inter-service communication over HTTP

The flexibility to demonstrate different technologies and patterns outweighed the operational complexity for a portfolio project.

---

## Architectural Thesis: *Longevity over Priority*

These services are designed to run on a **single, low-cost EC2 instance**.
The priority is not peak throughput but **staying alive under load, at minimal cost**.

Key ideas:

* **Heavily throttled & debounced APIs** — shaping traffic before it overwhelms the system.
* **Self-protecting services** — rate limits, concurrency caps, circuit breakers, and backpressure.
* **Resilience first** — prefer rejecting quickly over failing unpredictably.
* **Minimal cost footprint** — t3.micro instance with Docker Compose, Postgres, and Redis.

---

## Services

### Backend Services

* **observability-api** (port 8081): Standardized health checks, readiness, Prometheus metrics, and structured logging. Also provides shared observability handlers for other services.
* **feature-flags-api** (port 4000): Read-only boolean feature flags served from JSON files for local/prod environments.
* **outbox-api** (port 8080): Event delivery service using the Outbox pattern, integrating with feature flags for adaptive throttling.

### Frontend

* **portfolio-ui** (port 3000): Next.js React application showcasing the services and providing management interfaces.

### Database

* **postgres** (port 5432): PostgreSQL database for persistent data storage.

## Quick Start

All services are fully containerized and can be run individually or together:

### Individual Service Deployment

Each service can be deployed independently:

```bash
# Build and run observability-api
docker build -f packages/observability/Dockerfile -t observability-api:latest .
docker run --name observability-api -p 8081:8081 -d observability-api:latest

# Build and run feature-flags-api  
docker build -f apps/feature-flags-api/Dockerfile -t feature-flags-api:latest .
docker run --name feature-flags-api -p 4000:4000 -d feature-flags-api:latest

# Build and run outbox-api (requires postgres)
docker run --name postgres -e POSTGRES_DB=outbox -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15
docker build -f apps/outbox-api/Dockerfile -t outbox-api:latest .
docker run --name outbox-api -p 8080:8080 -e DB_HOST=host.docker.internal -e DB_PORT=5432 -e DB_USER=postgres -e DB_PASSWORD=password -e DB_NAME=outbox -e DB_SSLMODE=disable --add-host=host.docker.internal:host-gateway -d outbox-api:latest

# Build and run portfolio-ui
docker build -f ui/portfolio/Dockerfile -t portfolio-ui:latest .
docker run --name portfolio-ui -p 3000:3000 -d portfolio-ui:latest
```

### Full Stack Deployment

Use Docker Compose to run the entire stack:

```bash
docker-compose up --build -d
```

This will start all services with proper dependencies and health checks.

---

## Deployment Philosophy

* **Reverse proxy shaping:** Nginx/Caddy sits in front of APIs, shedding load before it hits application processes.
* **Application safeguards:** Concurrency limiter, timeouts, and circuit breakers.
* **Adaptive load shedding:** Services dynamically adjust worker batch sizes, rate limits, or concurrency caps based on health metrics.
* **Maintenance mode:** Services can be switched to low-power survival mode automatically or manually.
