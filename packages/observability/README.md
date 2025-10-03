# Observability API

The **Observability API** provides basic observability endpoints for services in this portfolio.  
It standardizes health checks, readiness checks, and Prometheus metrics exposure.

The observability-api provides a minimal baseline (health, readiness, HTTP counters, latency histograms, Go runtime metrics). As additional services are added, domain-specific metrics will be registered where they provide value.

## Dual Purpose

This package serves two purposes:

1. **Shared Library**: Provides handlers and middleware for other services to use
2. **Standalone Service**: Runs its own observability API server on port 8081

## Quick Start

### Prerequisites

- Go 1.25+
- Docker (recommended)

### Local Development

Start the service:

```bash
go run .
```

The service will start on port 8081 by default.

### Docker Deployment (Recommended)

The service is fully containerized and can be deployed as a standalone service:

1. **Build the Docker image**:

   ```bash
   # From the monorepo root
   docker build -f packages/observability/Dockerfile -t observability-api:latest .
   ```

2. **Run the observability-api container**:

   ```bash
   docker run --name observability-api \
     -p 8081:8081 \
     -d observability-api:latest
   ```

3. **Verify the service is running**:

   ```bash
   # Check container status
   docker ps
   
   # Test health endpoint
   curl http://localhost:8081/health
   
   # Test readiness endpoint
   curl http://localhost:8081/ready
   
   # Test metrics endpoint
   curl http://localhost:8081/metrics
   
   # View logs
   docker logs observability-api
   ```

4. **Stop the service**:

   ```bash
   docker stop observability-api
   docker rm observability-api
   ```

## Endpoints

### Health & Monitoring

- `GET /health` - Liveness probe (returns `{"status":"ok"}`)
- `GET /ready` - Readiness probe (returns `{"status":"ready"}`)
- `GET /metrics` - Prometheus metrics in text format

### API Documentation

- `GET /docs/index.html` - Swagger UI (OpenAPI docs)
- `GET /docs/doc.json` - OpenAPI specification in JSON format

## API Documentation

Interactive Swagger documentation is available when the service is running:

- **Swagger UI**: http://localhost:8081/docs/index.html
- **Raw JSON**: http://localhost:8081/docs/doc.json

### Example Usage

**Health check:**
```bash
curl http://localhost:8081/health
# Response: {"status":"ok"}
```

**Readiness check:**
```bash
curl http://localhost:8081/ready
# Response: {"status":"ready"}
```

**Prometheus metrics:**
```bash
curl http://localhost:8081/metrics
# Response: Prometheus-formatted metrics
```

## Metrics Provided

The service exposes the following Prometheus metrics:

### Go Runtime Metrics
- `go_gc_duration_seconds` - Garbage collection duration
- `go_goroutines` - Number of goroutines
- `go_memstats_*` - Memory statistics
- `go_threads` - Number of OS threads

### HTTP Metrics
- `http_requests_total` - Total HTTP requests by method, path, and status
- `http_request_duration_seconds` - HTTP request latency histogram

### Process Metrics
- `process_cpu_seconds_total` - CPU time consumed
- `process_memory_bytes` - Memory usage
- `process_open_fds` - Open file descriptors

## Configuration

The service runs with minimal configuration:

- **Port**: 8081 (hardcoded)
- **No external dependencies**: Self-contained service
- **No configuration files**: Uses defaults for all settings

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
