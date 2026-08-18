-- A-Radius / APB
-- Initial database migration
--
-- This migration is intentionally self-contained.
-- Schema files remain useful as documentation/reference,
-- but a fresh database must be bootstrappable from migrations.


CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS apb;

CREATE TABLE IF NOT EXISTS apb.schema_migrations (
    version VARCHAR(100) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE IF NOT EXISTS apb.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(150),
    phone VARCHAR(30),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_status_check
        CHECK (status IN ('active', 'inactive', 'suspended', 'locked'))
);

CREATE TABLE IF NOT EXISTS apb.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS apb.user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES apb.users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES apb.roles(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS apb.permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    permission_key VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS apb.role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY (role_id)
        REFERENCES apb.roles(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES apb.permissions(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS apb.security_knowledge_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    knowledge_key VARCHAR(150) NOT NULL,
    version VARCHAR(50) NOT NULL,
    content_hash VARCHAR(128) NOT NULL,
    source VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    learned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (knowledge_key, version),

    CONSTRAINT security_knowledge_status_check
        CHECK (status IN ('draft', 'active', 'deprecated', 'revoked'))
);

CREATE INDEX IF NOT EXISTS idx_user_roles_role
    ON apb.user_roles(role_id);

CREATE INDEX IF NOT EXISTS idx_role_permissions_permission
    ON apb.role_permissions(permission_id);

CREATE INDEX IF NOT EXISTS idx_users_status
    ON apb.users(status);

CREATE INDEX IF NOT EXISTS idx_users_email
    ON apb.users(email);

CREATE INDEX IF NOT EXISTS idx_audit_logs_actor
    ON apb.audit_logs(actor_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created
    ON apb.audit_logs(created_at);

CREATE INDEX IF NOT EXISTS idx_audit_logs_resource
    ON apb.audit_logs(resource_type, resource_id);

CREATE INDEX IF NOT EXISTS idx_security_knowledge_key
    ON apb.security_knowledge_versions(knowledge_key);

INSERT INTO apb.schema_migrations(version)
VALUES ('0001_init')
ON CONFLICT (version) DO NOTHING;
