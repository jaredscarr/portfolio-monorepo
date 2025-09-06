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

## Usage

Start the service:

```bash
go run .
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
