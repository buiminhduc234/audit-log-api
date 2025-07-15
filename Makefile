.PHONY: build run test clean migrate migrate-down

# Go related variables
BINARY_NAME=audit-log-api
MAIN_PACKAGE=./cmd/api

# Docker related variables
DOCKER_COMPOSE=docker-compose

# Build the application
build:
	@echo "Building audit-log-api..."
	@go build -o bin/audit-log-api ./cmd/api

# Run the application
run:
	@go run ./cmd/api

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@rm -rf bin/

# Start Docker services
docker-up:
	@echo "Starting Docker services..."
	$(DOCKER_COMPOSE) up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	$(DOCKER_COMPOSE) down

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	@sql-migrate down

migrate-up:
	@sql-migrate up -config=dbconfig.yml
	
# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...
	go fmt ./...

# Build and run the application
dev: docker-up deps build run

# Generate token
generate-token:
	@go run scripts/generate_token.go -user=11111111-1111-1111-1111-111111111111 -roles=user,auditor -tenant=11111111-1111-1111-1111-111111111111 -exp=1