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

## Architectural Thesis: *Longevity over Priority*

These services are designed to run on a **single, low-cost EC2 instance**.
The priority is not peak throughput but **staying alive under load, at minimal cost**.

Key ideas:

* **Heavily throttled & debounced APIs** — shaping traffic before it overwhelms the system.
* **Self-protecting services** — rate limits, concurrency caps, circuit breakers, and backpressure.
* **Resilience first** — prefer rejecting quickly over failing unpredictably.
* **Minimal cost footprint** — t3.micro instance with Docker Compose, Postgres, and Redis.

---

## Planned Services

* **metrics-api**: Standardized health checks, readiness, Prometheus metrics, and structured logging.
* **feature-flags-api**: Append-only feature flag service with Postgres + Redis, exposing CRUD, history, evaluation, and snapshot endpoints.
* **outbox-api**: Event delivery service using the Outbox pattern, integrating with feature flags for adaptive throttling.

---

## Deployment Philosophy

* **Reverse proxy shaping:** Nginx/Caddy sits in front of APIs, shedding load before it hits Go processes.
* **Application safeguards:** Concurrency limiter, timeouts, and circuit breakers.
* **Adaptive load shedding:** Services dynamically adjust worker batch sizes, rate limits, or concurrency caps based on health metrics.
* **Maintenance mode:** Services can be switched to low-power survival mode automatically or manually.
