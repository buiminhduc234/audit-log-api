.PHONY: build run test clean migrate migrate-down build-worker run-worker

# Go related variables
BINARY_NAME=audit-log-api
WORKER_BINARY_NAME=audit-log-worker
MAIN_PACKAGE=./cmd/api
WORKER_PACKAGE=./cmd/worker

# Docker related variables
DOCKER_COMPOSE=docker-compose

# Build the application
build:
	@echo "Building audit-log-api..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

# Build the worker
build-worker:
	@echo "Building audit-log-worker..."
	@go build -o bin/$(WORKER_BINARY_NAME) $(WORKER_PACKAGE)

# Build all
build-all: build build-worker

# Run the application
run:
	@go run $(MAIN_PACKAGE)

# Run the worker
run-worker:
	@go run $(WORKER_PACKAGE)

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
docker-clear:
	@echo "Stopping Docker services..."
	$(DOCKER_COMPOSE) down -v

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

# Build and run the worker
dev-worker: docker-up deps build-worker run-worker

# Generate token
generate-token:
	@go run scripts/generate_token.go -user=11111111-1111-1111-1111-111111111111 -roles=admin,user,auditor -tenant=11111111-1111-1111-1111-111111111111

swag:
	@echo '$(shell swag --version)'
	@swag init -g ./cmd/api/main.go --parseVendor true --exclude db,deployment,scripts,vendor
	@swagger2openapi ./docs/swagger.yaml -o ./docs/openapi.yaml
	@swagger2openapi ./docs/swagger.json -o ./docs/openapi.json