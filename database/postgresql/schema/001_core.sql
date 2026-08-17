-- A-Radius / APB
-- Core database foundation

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS apb;

CREATE TABLE IF NOT EXISTS apb.system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setting_key VARCHAR(150) NOT NULL UNIQUE,
    setting_value TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS apb.audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_actor
    ON apb.audit_logs(actor_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created
    ON apb.audit_logs(created_at);

CREATE INDEX IF NOT EXISTS idx_audit_logs_resource
    ON apb.audit_logs(resource_type, resource_id);
