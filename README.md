# Audit Log API

## Overview

A comprehensive audit logging API system designed to track and manage user actions across different applications with multi-tenant support. This system handles high-volume logging, provides advanced search and filtering capabilities, and ensures data integrity and security.

## Objective

The Audit Log API provides:
- **High-Performance Logging**: Handle 1000+ log entries per second with sub-100ms response times
- **Multi-Tenant Architecture**: Complete data isolation between tenants
- **Real-Time Streaming**: WebSocket-based live log monitoring
- **Advanced Search**: Full-text search and filtering capabilities via OpenSearch
- **Data Lifecycle Management**: Automated archival, cleanup, and retention policies
- **Enterprise Security**: JWT authentication, role-based access control, and data encryption

## Tech Stack

### Core Technologies
- **Language**: Go 1.21+
- **Web Framework**: Gin (HTTP router and middleware)
- **API Documentation**: Swagger/OpenAPI 3.0

### Data Storage
- **Primary Database**: PostgreSQL 15+ with TimescaleDB extension (optimized for time-series data)
- **Search Engine**: OpenSearch (advanced search and full-text search)
- **PubSub**: Redis ()
- **Archive Storage**: AWS S3 (long-term log storage)

### Message Queue & Workers
- **Queue System**: AWS SQS (background task processing)
- **Worker Services**: 
  - Index Worker (OpenSearch indexing)
  - Archive Worker (S3 archival)
  - Cleanup Worker (data retention)

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Local Development**: LocalStack (AWS services simulation)
- **Database Migration**: sql-migrate
- **Configuration**: Environment-based configuration

## Prerequisites

Before running the application locally, ensure you have:

- **Go 1.21+** installed
- **Docker** and **Docker Compose** installed
- **Make** utility (for running Makefile commands)

## Local Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd audit-log-api
```

### 2. Start Infrastructure Services

Start all required services using Docker Compose:

```bash
# Start PostgreSQL, OpenSearch, Redis, and LocalStack
make docker-up

# Wait for services to be ready
```

### 3. Initialize AWS Resources

Set up SQS queues and S3 buckets in LocalStack:

```bash
# Initialize SQS queues and S3 buckets
make init-localstack
```

### 5. Database Setup

Run database migrations:

```bash
# Create database schema and seed initial data
make migrate-up
```

### 6. Build the Application

```bash
# Build all services
make build-all
```

## Running the Application

### Start All Services

```bash
# Start the main API server
make run-api

# In separate terminals, start the workers:
make run-index-worker    # OpenSearch indexing
make run-archive-worker  # S3 archival
make run-cleanup-worker  # Data cleanup
```

### Verify Installation

1. **API Health Check**:
   ```bash
   curl http://localhost:10000/health
   ```

2. **API Documentation**:
   Open http://localhost:10000/swagger/index.html in your browser

3. **Generate Test Token**:
   ```bash
   make generate-token
   ```

4. **Test API Endpoints**:
   Import this [Postman collection](docs/AuditLogAPI.postman_collection.json) for testing

## Project Structure

```
audit-log-api/
├── cmd/                    # Application entry points
│   ├── api/               # Main API server
│   ├── archive_worker/    # S3 archive worker
│   ├── cleanup_worker/    # Data cleanup worker
│   └── index_worker/      # OpenSearch index worker
├── internal/              # Internal application code
│   ├── api/              # HTTP handlers and routes
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models
│   ├── middleware/       # HTTP middleware
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── worker/           # Background workers
├── pkg/                   # Public packages
├── scripts/              # Database migrations and utilities
├── docs/                 # API documentation
└── docker-compose.yml    # Local development services
```

## Documentation

For detailed information about the system, refer to:

- **[AUDIT_LOG_FLOW_DIAGRAMS.md](AUDIT_LOG_FLOW_DIAGRAMS.md)** - System architecture and data flow
- **[QUEUE_ARCHITECTURE.md](QUEUE_ARCHITECTURE.md)** - queue architecture
- **API Documentation**: http://localhost:10000/swagger/index.html (when running)

## Features

- ✅ **Multi-tenant Architecture** with complete data isolation
- ✅ **High-Performance API** (1000+ requests/second)
- ✅ **Real-time WebSocket Streaming** for live log monitoring
- ✅ **Advanced Search** with OpenSearch integration
- ✅ **Data Lifecycle** (archival, cleanup, retention)
- ✅ **JWT Authentication** with role-based access control
- ✅ **AWS Integration** (SQS, S3) with LocalStack support
- ✅ **Comprehensive API Documentation** with OpenAPI/Swagger
- ✅ **Database Read/Write Separation** for optimal performance
- ✅ **Background Workers** for async processing
