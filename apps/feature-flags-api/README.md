# Feature Flags API

A lightweight HTTP service for serving **read-only boolean feature flags** for `local` and `prod` environments.  
Backed by static JSON files committed to the monorepo.

## Purpose

This service provides a central source of truth for environment-specific feature flags.

## How It Works

- Flags are stored as JSON files in the monorepo under `flags/`
  - Example: `flags/local.json`, `flags/prod.json`
- On startup, the service loads both flag files into memory
- Flags can be hot-reloaded via an admin endpoint without restarting the service
- Swagger docs are available at: [http://localhost:4000/swagger/index.html](http://localhost:4000/swagger/index.html)

## Quick Start

### Prerequisites

- Go 1.25+
- Docker (recommended)

### Local Development

Start the service:

```bash
go run .
```

The service will start on port 4000 by default.

### Docker Deployment (Recommended)

The service is fully containerized and can be deployed as a standalone service:

1. **Build the Docker image**:

   ```bash
   # From the monorepo root
   docker build -f apps/feature-flags-api/Dockerfile -t feature-flags-api:latest .
   ```

2. **Run the feature-flags-api container**:

   ```bash
   docker run --name feature-flags-api \
     -p 4000:4000 \
     -d feature-flags-api:latest
   ```

3. **Verify the service is running**:

   ```bash
   # Check container status
   docker ps
   
   # Test health endpoint
   curl http://localhost:4000/health
   
   # Test flags endpoint
   curl "http://localhost:4000/flags?env=local"
   
   # View logs
   docker logs feature-flags-api
   ```

4. **Stop the service**:

   ```bash
   docker stop feature-flags-api
   docker rm feature-flags-api
   ```

## API Endpoints

### Health & Monitoring

- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check endpoint  
- `GET /metrics` - Prometheus metrics endpoint

### Feature Flags

- `GET /flags?env={local|prod}` - Get all flags for environment
- `GET /flags/:key?env={local|prod}` - Get specific flag by key

### Administration

- `POST /admin/reload` - Reload flags from disk
- `PUT /admin/flags/:key` - Update a flag value

## API Documentation

Interactive Swagger documentation is available when the service is running:

- **Swagger UI**: http://localhost:4000/swagger/index.html
- **Raw JSON**: http://localhost:4000/swagger/doc.json

### Example Usage

**Get all local flags:**
```bash
curl "http://localhost:4000/flags?env=local"
```

**Get specific flag:**
```bash
curl "http://localhost:4000/flags/simulation_mode_enabled?env=local"
```

**Reload flags:**
```bash
curl -X POST http://localhost:4000/admin/reload
```

## Configuration

The service uses static JSON files for configuration. Flag files are located in the `flags/` directory:

- `flags/local.json` - Local environment flags
- `flags/prod.json` - Production environment flags

### Flag File Format

```json
{
  "simulation_mode_enabled": true,
  "advanced_debugging_enabled": true,
  "circuit_breaker_demo_mode": false,
  "disable_publishing": false,
  "force_webhook_failures": false,
  "metrics_enabled": false,
  "partial_failure_mode": false,
  "simulate_network_delays": false
}
```

## Testing

This project includes comprehensive unit tests for all handlers and admin functions.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./handlers -v

# Run tests with coverage
go test ./... -cover

# Run tests with coverage and generate HTML report
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run benchmarks with memory allocation info
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkGetFlags -benchmem

# Run benchmarks for longer (more accurate results)
go test -bench=. -benchtime=10s
```

## Update Swagger docs

```bash
swag init -g main.go -o docs/ -ot json
