APP_NAME=metrics-api

.PHONY: run swag tidy

# Run the service
run:
	go run main.go

# Regenerate Swagger docs
swag:
	swag init --parseDependency --parseInternal

# Clean up dependencies
tidy:
	go mod tidy
