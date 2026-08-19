CREATE SCHEMA IF NOT EXISTS security;

CREATE TYPE security.knowledge_status AS ENUM (
    'DISCOVERED', 'ANALYZING', 'VALIDATED', 'REVIEW_REQUIRED',
    'APPROVED', 'STAGED', 'ACTIVE', 'SUPERSEDED', 'ARCHIVED'
);

CREATE TYPE security.approval_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED');

CREATE TABLE security.security_knowledge_versions (
    id BIGSERIAL PRIMARY KEY,
    version TEXT NOT NULL UNIQUE,
    major INTEGER NOT NULL CHECK (major >= 0),
    minor INTEGER NOT NULL CHECK (minor >= 0),
    patch INTEGER NOT NULL CHECK (patch >= 0),
    status security.knowledge_status NOT NULL DEFAULT 'DISCOVERED',
    checksum TEXT NOT NULL,
    source_count INTEGER NOT NULL DEFAULT 0 CHECK (source_count >= 0),
    finding_count INTEGER NOT NULL DEFAULT 0 CHECK (finding_count >= 0),
    confidence_score NUMERIC(5,2) CHECK (confidence_score >= 0 AND confidence_score <= 100),
    environment TEXT NOT NULL DEFAULT 'ai-knowledge-db' CHECK (environment IN ('ai-knowledge-db', 'analysis-cache', 'staging', 'production')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    discovered_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    staged_at TIMESTAMPTZ,
    activated_at TIMESTAMPTZ,
    superseded_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    production_changed BOOLEAN NOT NULL DEFAULT false,
    CHECK (production_changed = false OR status IN ('STAGED', 'ACTIVE', 'SUPERSEDED', 'ARCHIVED'))
);

CREATE UNIQUE INDEX security_one_active_knowledge_version
    ON security.security_knowledge_versions (status)
    WHERE status = 'ACTIVE';

CREATE TABLE security.security_intelligence (
    id BIGSERIAL PRIMARY KEY,
    knowledge_version_id BIGINT NOT NULL REFERENCES security.security_knowledge_versions(id) ON DELETE CASCADE,
    intelligence_key TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL,
    severity TEXT NOT NULL CHECK (severity IN ('CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'INFO')),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    source TEXT NOT NULL,
    evidence JSONB NOT NULL DEFAULT '{}'::jsonb,
    confidence NUMERIC(5,2) CHECK (confidence >= 0 AND confidence <= 100),
    discovered_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE security.security_rules (
    id BIGSERIAL PRIMARY KEY,
    knowledge_version_id BIGINT NOT NULL REFERENCES security.security_knowledge_versions(id) ON DELETE CASCADE,
    rule_code TEXT NOT NULL,
    category TEXT NOT NULL,
    severity TEXT NOT NULL CHECK (severity IN ('CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'INFO')),
    rule_definition JSONB NOT NULL,
    recommendation TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    UNIQUE (knowledge_version_id, rule_code)
);

CREATE TABLE security.security_application_mapping (
    id BIGSERIAL PRIMARY KEY,
    knowledge_version_id BIGINT NOT NULL REFERENCES security.security_knowledge_versions(id) ON DELETE CASCADE,
    module TEXT NOT NULL,
    component TEXT NOT NULL,
    risk TEXT NOT NULL CHECK (risk IN ('CRITICAL', 'HIGH', 'MEDIUM', 'LOW', 'INFO')),
    affected BOOLEAN NOT NULL DEFAULT true,
    UNIQUE (knowledge_version_id, module, component)
);

CREATE TABLE security.security_knowledge_approvals (
    id BIGSERIAL PRIMARY KEY,
    knowledge_version_id BIGINT NOT NULL REFERENCES security.security_knowledge_versions(id) ON DELETE CASCADE,
    action TEXT NOT NULL CHECK (action IN ('APPROVE', 'STAGE', 'ACTIVATE', 'ROLLBACK', 'ARCHIVE')),
    requested_by TEXT NOT NULL,
    approved_by TEXT,
    status security.approval_status NOT NULL DEFAULT 'PENDING',
    reason TEXT NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at TIMESTAMPTZ,
    CHECK (status <> 'APPROVED' OR approved_by IS NOT NULL)
);

CREATE TABLE security.security_knowledge_audit (
    id BIGSERIAL PRIMARY KEY,
    knowledge_version_id BIGINT REFERENCES security.security_knowledge_versions(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    actor_type TEXT NOT NULL CHECK (actor_type IN ('AI', 'DEVELOPER', 'ADMINISTRATOR', 'SYSTEM')),
    actor_id TEXT NOT NULL,
    old_value JSONB NOT NULL DEFAULT '{}'::jsonb,
    new_value JSONB NOT NULL DEFAULT '{}'::jsonb,
    correlation_id TEXT,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX security_intelligence_version_idx ON security.security_intelligence (knowledge_version_id);
CREATE INDEX security_rules_version_idx ON security.security_rules (knowledge_version_id);
CREATE INDEX security_mapping_version_idx ON security.security_application_mapping (knowledge_version_id);
CREATE INDEX security_approvals_version_idx ON security.security_knowledge_approvals (knowledge_version_id);
CREATE INDEX security_audit_version_time_idx ON security.security_knowledge_audit (knowledge_version_id, occurred_at DESC);
