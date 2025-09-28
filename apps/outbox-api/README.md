# Outbox Pattern Publishing Service

A reliable event publishing service implementing the Outbox Pattern for guaranteed event delivery in distributed systems.

## Overview

The Outbox Pattern Publishing Service provides a robust solution for publishing events from your application with guaranteed delivery. It implements the Outbox Pattern to ensure that events are reliably published even in the face of failures.

## Features

- **Reliable Event Publishing**: Implements the Outbox Pattern for guaranteed event delivery
- **Batch Processing**: Efficiently processes events in configurable batches
- **Circuit Breaker**: Built-in circuit breaker for external service protection
- **Retry Logic**: Configurable retry attempts with exponential backoff
- **Observability**: Integrated health checks, metrics, and monitoring
- **Self-Protecting**: Rate limiting, adaptive load shedding, and backpressure handling

## Architecture

The service follows a "self-protecting service" deployment strategy focused on longevity over peak throughput:

- **Standalone Service**: Independent HTTP API service
- **Shared Database**: PostgreSQL for event storage and state management
- **Batch Publishing**: Configurable batch processing for efficiency
- **Circuit Breaker**: Protection against external service failures
- **Observability**: Health checks, metrics, and monitoring endpoints

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker (recommended)

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/jared-scarr/portfolio-monorepo.git
   cd portfolio-monorepo
   ```

2. **Set up the workspace**:

   ```bash
   go work init
   go work use ./apps/outbox-api
   go work use ./packages/observability
   ```

3. **Install dependencies**:

   ```bash
   cd apps/outbox-api
   go mod tidy
   ```

4. **Configure the service**:

   ```bash
   cp config.json.example config.json
   # Edit config.json with your database settings
   ```

5. **Set up the database**:

   ```bash
   # Create PostgreSQL database
   createdb outbox
   
   # Or use Docker
   docker run --name postgres-outbox -e POSTGRES_PASSWORD=password -e POSTGRES_DB=outbox -p 5432:5432 -d postgres:15
   ```

6. **Run the service**:

   ```bash
   go run .
   ```

The service will start on port 8080 by default.

### Docker Deployment (Recommended)

The service is fully containerized and can be run using Docker Compose for easy setup:

1. **Start the complete stack**:

   ```bash
   cd apps/outbox-api
   docker-compose up --build -d
   ```

2. **Verify the services are running**:

   ```bash
   docker-compose ps
   ```

3. **Check the health endpoint**:

   ```bash
   curl http://localhost:8080/health
   ```

4. **View logs**:

   ```bash
   docker-compose logs -f outbox-api
   ```

5. **Stop the services**:

   ```bash
   docker-compose down
   ```

#### Docker Services

The `docker-compose.yml` includes:

- **outbox-api**: The main application service
- **postgres**: PostgreSQL 15 database
- **Automatic health checks**: Services wait for dependencies to be healthy
- **Environment configuration**: Database connection automatically configured

#### Building the Docker Image

To build the Docker image manually:

```bash
# From the monorepo root
docker build -f apps/outbox-api/Dockerfile -t outbox-api:latest .
```

#### Environment Configuration

See [ENVIRONMENT.md](ENVIRONMENT.md) for detailed environment variable configuration.

## Configuration

The service uses environment variables for configuration. Create a `.env` file in the project root:

```bash
# Database Configuration
DB_USER=postgres
DB_PASSWORD=password
DB_HOST=localhost
DB_NAME=outbox

# Server Configuration
PORT=8080

# Webhook Configuration
WEBHOOK_URL=http://localhost:3000/webhook

# Optional overrides
BATCH_SIZE=10
```



### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `outbox` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `WEBHOOK_URL` | Webhook endpoint URL | `http://localhost:3000/webhook` |
| `BATCH_SIZE` | Batch size for publishing | `10` |

### Production Security

For production deployments, use environment variables for sensitive data:

```bash
export DB_PASSWORD="your_secure_password"
export DB_USER="your_db_user"
export DB_HOST="your_db_host"
export WEBHOOK_URL="https://your-secure-webhook.com/endpoint"
```

This prevents sensitive data from being stored in configuration files that might be committed to version control.

## API Endpoints

### Health & Monitoring

- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check endpoint  
- `GET /metrics` - Prometheus metrics endpoint

### Event Management

- `POST /api/v1/events` - Create a new event
- `GET /api/v1/events` - List events (with pagination and filtering)
- `GET /api/v1/events/:id` - Get event by ID
- `POST /api/v1/events/:id/retry` - Retry a failed event
- `DELETE /api/v1/events/:id` - Delete an event

### Administration

- `POST /admin/publish` - Manually trigger event publishing
- `GET /admin/stats` - Get service statistics

## API Documentation

Interactive Swagger documentation is available when the service is running:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Raw JSON**: http://localhost:8080/swagger/doc.json

### Updating API Documentation

When you modify API endpoints or add new ones, regenerate the Swagger documentation:

```bash
# Install swag CLI tool (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swag init -g main.go -o docs/ -ot json

# Restart the server to serve updated docs
go run .
```

The documentation will automatically reflect changes to endpoint handlers and their Swagger annotations.

## Usage Examples

### Creating an Event

**Linux/macOS:**

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user.created",
    "source": "user-service",
    "data": {
      "user_id": "123",
      "email": "user@example.com",
      "name": "John Doe"
    },
    "metadata": {
      "version": "1.0",
      "correlation_id": "abc-123"
    }
  }'
```

**Windows PowerShell:**

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/events" -Method POST -ContentType "application/json" -Body '{"type": "user.created", "source": "user-service", "data": {"user_id": "123", "email": "user@example.com", "name": "John Doe"}, "metadata": {"version": "1.0", "correlation_id": "abc-123"}}'
```

### Listing Events

**Linux/macOS:**

```bash
# List all events
curl http://localhost:8080/api/v1/events

# List pending events only
curl http://localhost:8080/api/v1/events?status=pending

# Paginated results
curl http://localhost:8080/api/v1/events?page=1&limit=10
```

**Windows PowerShell:**

```powershell
# List all events
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/events"

# List pending events only
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/events?status=pending"

# Paginated results
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/events?page=1&limit=10"
```

### Getting Statistics

**Linux/macOS:**

```bash
curl http://localhost:8080/admin/stats
```

**Windows PowerShell:**

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/admin/stats"
```

## Development

### Running Tests

```bash
go test ./...
```

### Running with Coverage

```bash
go test -cover ./...
```

### Building

```bash
go build -o outbox-api .
```

## Deployment

The service is designed for deployment on AWS EC2 `t3.micro` instances with:

- **Reverse Proxy**: Caddy or Nginx for SSL termination and load balancing
- **Database**: PostgreSQL (can be RDS or self-managed)
- **Monitoring**: CloudWatch Agent for metrics and logs
- **Containerization**: Docker for consistent deployments

## Self-Protection Features

The service implements several self-protection mechanisms:

1. **Circuit Breaker**: Prevents cascading failures from external services
2. **Rate Limiting**: Protects against overwhelming external endpoints
3. **Exponential Backoff**: Intelligent retry delays for failed requests
4. **Adaptive Batching**: Dynamic batch sizing based on system load
5. **Backpressure**: Graceful degradation under high load

## Future Enhancements

- **AWS Integration**: Migration to SQS/EventBridge for cloud-native event publishing
- **Advanced Filtering**: Event filtering and routing capabilities
- **Dead Letter Queue**: Handling of permanently failed events
- **Event Schema Validation**: JSON schema validation for events
- **Multi-tenant Support**: Isolated event streams per tenant

## License

This project is part of the portfolio-monorepo and follows the same license terms.
