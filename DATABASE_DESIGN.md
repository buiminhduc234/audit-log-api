# Database Design: Audit Log API

This document describes the **database schema and design decisions** for the Audit Log API project.  
The database is implemented using **PostgreSQL with TimescaleDB extension**, optimized for **multi-tenant, time-series audit log data**.

---

## Schema Overview

### `tenants` table
Manages tenant metadata for multi-tenant support.

| Column        | Type         | Description                         |
|----------------|--------------|-------------------------------------|
| `id`           | UUID         | Primary key, auto-generated         |
| `name`         | TEXT         | Tenant name                         |
| `rate_limit`   | INTEGER      | Requests per second allowed         |
| `created_at`   | TIMESTAMPTZ  | Row creation timestamp              |
| `updated_at`   | TIMESTAMPTZ  | Row update timestamp                |

---

### `audit_logs` table
Stores audit log entries, optimized for time-series workloads using TimescaleDB hypertable.

| Column          | Type         | Description                          |
|------------------|--------------|--------------------------------------|
| `id`            | UUID         | Primary key (with `timestamp`)      |
| `tenant_id`     | UUID         | References `tenants(id)`            |
| `message`       | TEXT         | Human-readable log message          |
| `user_id`       | TEXT         | ID of the user performing action    |
| `session_id`    | TEXT         | Session identifier                  |
| `ip_address`    | TEXT         | IP address of the client            |
| `user_agent`    | TEXT         | User agent string                   |
| `action`        | TEXT         | Action type (CREATE, UPDATE, etc.) |
| `resource_type` | TEXT         | Resource type affected              |
| `resource_id`   | TEXT         | Resource ID affected                |
| `severity`      | TEXT         | Severity level (INFO, ERROR, etc.) |
| `before_state`  | JSONB        | State of resource before change     |
| `after_state`   | JSONB        | State of resource after change      |
| `metadata`      | JSONB        | Additional structured metadata      |
| `timestamp`     | TIMESTAMPTZ  | Logical event timestamp             |
| `created_at`    | TIMESTAMPTZ  | Row creation timestamp              |
| `updated_at`    | TIMESTAMPTZ  | Row update timestamp                |

Primary Key: (`id`, `timestamp`) — supports time-series optimization and uniqueness.

---

## Multi-Tenancy

- Every `audit_logs` row is associated with a `tenant_id`, ensuring tenant isolation.  
- Foreign key with `ON DELETE CASCADE` ensures cleanup when a tenant is removed.

---

## Indexing Strategy

To optimize query performance:
- BRIN index on `timestamp` for efficient range queries.
- Partial index on `resource_type` where not null.
- Aggregation-friendly indexes on (`tenant_id`, `timestamp`, `action`), (`tenant_id`, `timestamp`, `severity`).

---

## TimescaleDB Features

### Hypertable
- `audit_logs` is converted into a **hypertable**, partitioned by `timestamp` in **1-day chunks**.

### Compression
- Chunks older than **7 days** are automatically compressed.
- Segments by `tenant_id, action, severity, resource_type` for efficient storage and decompression.

---

## Continuous Aggregates

### `audit_logs_hourly_stats`
A materialized view using TimescaleDB continuous aggregates:
- Aggregates log counts hourly per tenant, action, severity, and resource type.
- Covers up to 1 month of data.
- Automatically refreshed every hour.

This enables fast dashboard queries without scanning raw logs.
