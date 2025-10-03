# Multi-stage Dockerfile for feature-flags-api
# Stage 1: Build stage
FROM golang:1.25.0-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Install swag for generating docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set working directory
WORKDIR /app

# Copy go.mod files for dependency resolution
COPY apps/feature-flags-api/go.mod apps/feature-flags-api/go.sum ./
COPY packages/observability/go.mod packages/observability/go.sum ../packages/observability/

# Download dependencies
RUN go mod download

# Copy source code
COPY apps/feature-flags-api/ ./
COPY packages/observability/ ../packages/observability/

# Generate swagger docs
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o feature-flags-api .

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary and flag files from builder stage
COPY --from=builder /app/feature-flags-api .
COPY --from=builder /app/flags ./flags

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

# Run the application
CMD ["./feature-flags-api"]