# Environment Configuration for outbox-api

## Environment Variables

The outbox-api service uses environment variables for configuration. Here are the available options:

### Database Configuration
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database username (default: postgres)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: outbox)
- `DB_SSLMODE` - SSL mode (default: disable)

### Server Configuration
- `PORT` - Server port (default: 8080)
- `READ_TIMEOUT` - Read timeout (default: 30s)
- `WRITE_TIMEOUT` - Write timeout (default: 30s)

### Webhook Configuration
- `WEBHOOK_URL` - Webhook endpoint URL (default: http://localhost:3000/webhook)

### Publishing Configuration (Optional)
- `BATCH_SIZE` - Batch size for publishing (default: 10)
- `BATCH_TIMEOUT` - Batch timeout (default: 5s)
- `RETRY_ATTEMPTS` - Number of retry attempts (default: 3)
- `RETRY_DELAY` - Initial retry delay (default: 1s)
- `MAX_RETRY_DELAY` - Maximum retry delay (default: 30s)

### Circuit Breaker Configuration (Optional)
- `CIRCUIT_MAX_REQUESTS` - Maximum requests before circuit opens (default: 5)
- `CIRCUIT_INTERVAL` - Circuit breaker interval (default: 10s)
- `CIRCUIT_TIMEOUT` - Circuit breaker timeout (default: 5s)

### Development Settings
- `GIN_MODE` - Gin mode (debug/release, default: debug)

## Docker Compose Environment

When using docker-compose, the following environment variables are automatically set:

```yaml
environment:
  DB_HOST: postgres
  DB_PORT: 5432
  DB_USER: postgres
  DB_PASSWORD: password
  DB_NAME: outbox
  DB_SSLMODE: disable
  PORT: 8080
  WEBHOOK_URL: http://host.docker.internal:3000/webhook
```

## Local Development

For local development, create a `.env` file in the `apps/outbox-api` directory:

```bash
# Copy and modify as needed
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=outbox
DB_SSLMODE=disable
PORT=8080
WEBHOOK_URL=http://localhost:3000/webhook
```

## Production Deployment

For production, use environment variables or a secrets management system:

```bash
export DB_PASSWORD="your_secure_password"
export DB_USER="your_db_user"
export DB_HOST="your_db_host"
export WEBHOOK_URL="https://your-secure-webhook.com/endpoint"
```
