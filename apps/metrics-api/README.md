# Metrics API

The **Metrics API** provides basic observability endpoints for services in this portfolio.  
It standardizes health checks, readiness checks, and Prometheus metrics exposure.

## Endpoints

- `GET /health` → Liveness probe  
- `GET /ready` → Readiness probe  
- `GET /metrics` → Prometheus metrics (text format)  
- `GET /docs/index.html` → Swagger UI (OpenAPI docs)  
- `GET /openapi.json` → OpenAPI spec  

## Running Locally

```bash
make run

## Updating Swagger Docs

```bash
make swag

## Dependency Management
```bash
make tidy

