# Metrics API

The **Metrics API** provides basic observability endpoints for services in this portfolio.  
It standardizes health checks, readiness checks, and Prometheus metrics exposure.

The metrics-api provides a minimal baseline (health, readiness, HTTP counters, latency histograms, Go runtime metrics). As additional services are added, domain-specific metrics will be registered where they provide value.

## Endpoints

- `GET /health` → Liveness probe  
- `GET /ready` → Readiness probe  
- `GET /metrics` → Prometheus metrics (text format)  
- `GET /docs/index.html` → Swagger UI (OpenAPI docs)  
- `GET /openapi.json` → OpenAPI spec  

## Running Locally

```bash
go run .
```

## Testing

This project includes comprehensive unit tests for all handlers and middleware.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./internal/handlers -v

# Run tests with coverage
go test ./... -cover

# Run tests with coverage and generate HTML report
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run benchmarks with memory allocation info
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkHealth -benchmem

# Run benchmarks for longer (more accurate results)
go test -bench=. -benchtime=10s
```

## Update Swagger docs

```bash
swag init -g main.go -o docs/ -ot json
