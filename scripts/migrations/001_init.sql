-- +migrate Up
-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create tenants table for multi-tenant support
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    rate_limit INTEGER NOT NULL DEFAULT 1000,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create audit_logs table optimized for time-series data
CREATE TABLE IF NOT EXISTS audit_logs (
    -- Primary identification
    id UUID NOT NULL,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    
    -- User information
    user_id TEXT,
    session_id TEXT,
    ip_address TEXT,
    user_agent TEXT,
    
    -- Action details
    action TEXT NOT NULL,
    resource_type TEXT,
    resource_id TEXT,
    severity TEXT NOT NULL DEFAULT 'INFO',
    
    -- State changes
    before_state JSONB,
    after_state JSONB,
    
    -- Additional data
    metadata JSONB,
    
    -- Timestamps
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Primary key using both id and timestamp for TimescaleDB
    CONSTRAINT pk_audit_logs PRIMARY KEY (id, timestamp)
);

-- Create indexes for audit_logs
CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id, timestamp DESC);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(tenant_id, user_id, timestamp DESC);
CREATE INDEX idx_audit_logs_session ON audit_logs(tenant_id, session_id, timestamp DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(tenant_id, action, timestamp DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(tenant_id, resource_type, resource_id, timestamp DESC);
CREATE INDEX idx_audit_logs_severity ON audit_logs(tenant_id, severity, timestamp DESC);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX idx_audit_logs_brin_timestamp ON audit_logs USING BRIN (timestamp);

-- Convert audit_logs to hypertable
SELECT create_hypertable('audit_logs', 'timestamp', chunk_time_interval => INTERVAL '1 day', if_not_exists => TRUE);

-- Create a function to clean up old audit logs
-- CREATE OR REPLACE FUNCTION delete_old_audit_logs(retention_days INTEGER)
-- RETURNS void AS $$
-- BEGIN
--     DELETE FROM audit_logs
--     WHERE timestamp < NOW() - (retention_days || ' days')::INTERVAL;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- Add a default retention policy (90 days)
-- SELECT add_retention_policy('audit_logs', INTERVAL '90 days', if_not_exists => TRUE);

-- +migrate Down

-- Drop retention policy
SELECT remove_retention_policy('audit_logs', if_exists => TRUE);

-- Drop the cleanup function
DROP FUNCTION IF EXISTS delete_old_audit_logs;

-- Drop tables in correct order due to foreign key constraints
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;

-- Drop TimescaleDB extension
DROP EXTENSION IF EXISTS timescaledb;
