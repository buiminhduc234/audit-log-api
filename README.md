# Audit Log API

A scalable audit logging service built with Go, PostgreSQL, and TimescaleDB.

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later
- Make (optional)

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/buiminhduc234/audit-log-api.git
   cd audit-log-api
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Start the development environment:
   ```bash
   make dev
   ```

   This command will:
   - Start the database containers
   - Install dependencies
   - Build the application
   - Run the application

## Build Instructions

### Using Make (Recommended)

The project includes a Makefile with common commands:

```bash
# Show available commands
make help

# Install dependencies
make deps

# Start database services
make docker-up

# Run database migrations
make migrate

# Build the application
make build

# Run the application
make run

# Run tests
make test

# Clean build artifacts
make clean

# Stop database services
make docker-down
```

### Manual Build

1. Install dependencies:
   ```bash
   go mod download
   go mod tidy
   ```

2. Build the application:
   ```bash
   go build -o bin/audit-log-api ./cmd/api
   ```

3. Run the application:
   ```bash
   ./bin/audit-log-api
   ```

## Database Setup

### Using Docker Compose (Recommended)

1. Start the database and pgAdmin:
   ```bash
   docker-compose up -d
   ```

   This will:
   - Start PostgreSQL with TimescaleDB (port 5432)
   - Start pgAdmin web interface (port 5050)
   - Run database migrations automatically
   - Configure TimescaleDB with optimized settings

2. Database connection details:
   ```
   Host: localhost
   Port: 5432
   Database: audit_log
   Username: postgres
   Password: postgres
   ```

3. PgAdmin access:
   ```
   URL: http://localhost:5050
   Email: admin@admin.com
   Password: admin
   ```

### Database Schema

The database includes two main tables:

1. `tenants` - Stores tenant information:
   - `id` (TEXT) - Primary key
   - `name` (TEXT) - Tenant name
   - `api_key` (TEXT) - Unique API key for authentication
   - `rate_limit` (INTEGER) - API rate limit per tenant
   - `created_at` (TIMESTAMP) - Creation timestamp
   - `updated_at` (TIMESTAMP) - Last update timestamp

2. `audit_logs` - Stores audit events:
   - `id` (TEXT) - Event ID
   - `tenant_id` (TEXT) - Foreign key to tenants
   - `user_id` (TEXT) - User who performed the action
   - `session_id` (TEXT) - Session identifier
   - `action` (TEXT) - Action performed
   - `resource` (TEXT) - Resource type
   - `resource_id` (TEXT) - Resource identifier
   - `ip_address` (TEXT) - IP address
   - `user_agent` (TEXT) - User agent string
   - `severity` (TEXT) - Event severity
   - `metadata` (JSONB) - Additional event data
   - `timestamp` (TIMESTAMP) - Event timestamp

## Project Structure

```
audit-log-api/
├── cmd/
│   └── api/              # Application entrypoint
├── internal/
│   ├── api/             # HTTP handlers
│   ├── config/          # Configuration
│   ├── domain/          # Domain models
│   ├── middleware/      # HTTP middleware
│   ├── repository/      # Data access layer
│   └── service/         # Business logic
├── pkg/
│   ├── logger/          # Logging package
│   └── utils/           # Shared utilities
├── scripts/
│   └── migrations/      # Database migrations
├── docker/              # Docker configuration
├── docker-compose.yml   # Docker services
├── Makefile            # Build automation
└── README.md           # Documentation
```

## Development

### Code Style

The project follows standard Go code style. Run the following before committing:

```bash
make lint
```

### Testing

Run the test suite:

```bash
make test
```

### Environment Variables

Copy `.env.example` to `.env` and adjust the values:

```env
# Server Configuration
SERVER_PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=audit_log

# JWT Configuration
JWT_SECRET_KEY=your-secret-key
JWT_EXPIRATION_HOURS=24
```

## Troubleshooting

1. If the containers don't start:
   ```bash
   # Check container logs
   docker-compose logs db
   docker-compose logs pgadmin
   ```

2. To reset the database:
   ```bash
   # Stop containers and remove volumes
   docker-compose down -v
   
   # Start fresh
   docker-compose up -d
   ```

3. Common issues:
   - Port conflicts: Make sure ports 5432 and 5050 are available
   - Permission issues: Check Docker daemon permissions
   - Connection refused: Wait a few seconds after container start