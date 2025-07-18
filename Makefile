.PHONY: build run test clean migrate migrate-down build-worker run-worker

DOCKER_COMPOSE=docker-compose

init-localstack:
	./scripts/init-localstack.sh

build:
	@echo "Building audit-log-api..."
	@go build -o bin/audit-log-api ./cmd/api/main.go

build-archive-worker:
	@echo "Building archive-worker..."
	@go build -o bin/archive_worker ./cmd/archive_worker

build-cleanup-worker:
	@echo "Building cleanup-worker..."
	@go build -o bin/cleanup_worker ./cmd/cleanup_worker

build-index-worker:
	@echo "Building index-worker..."
	@go build -o bin/index_worker ./cmd/index_worker

build-all: build build-index-worker build-archive-worker build-cleanup-worker

run-api:
	@go run ./cmd/api/main.go

run-index-worker:
	@go run ./cmd/index_worker

run-archive-worker:
	@go run ./cmd/archive_worker

run-cleanup-worker:
	@go run ./cmd/cleanup_worker

test:
	@go test -v ./...

clean:
	@rm -rf bin/

docker-up:
	@echo "Starting Docker services..."
	$(DOCKER_COMPOSE) up -d

docker-clear:
	@echo "Stopping Docker services..."
	$(DOCKER_COMPOSE) down -v

migrate-down:
	@echo "Rolling back database migrations..."
	@sql-migrate down

migrate-up:
	@sql-migrate up -config=dbconfig.yml
	
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

lint:
	@echo "Running linter..."
	go vet ./...
	go fmt ./...

# Generate token
generate-token:
	@go run ./cmd/app/generate_token.go -user=11111111-1111-1111-1111-111111111111 -roles=admin,user,auditor -tenant=11111111-1111-1111-1111-111111111111

swag:
	@echo '$(shell swag --version)'
	@swag init -g ./cmd/api/main.go --parseVendor true --exclude db,deployment,scripts,vendor
	@swagger2openapi ./docs/swagger.yaml -o ./docs/openapi.yaml
	@swagger2openapi ./docs/swagger.json -o ./docs/openapi.json