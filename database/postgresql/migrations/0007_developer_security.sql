-- Findings
CREATE TABLE IF NOT EXISTS security_findings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    finding_key     TEXT NOT NULL UNIQUE,
    title           TEXT NOT NULL,
    severity        TEXT NOT NULL CHECK (severity IN ('CRITICAL','HIGH','MEDIUM','LOW','INFO')),
    module          TEXT,
    status          TEXT NOT NULL DEFAULT 'open'
                    CHECK (status IN ('open','in_review','approved','ignored','fixed')),
    description     TEXT,
    recommendation  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Knowledge Versions
CREATE TABLE IF NOT EXISTS security_knowledge_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version         TEXT NOT NULL UNIQUE,
    status          TEXT NOT NULL CHECK (status IN ('new','review','active','archived')),
    findings_count  INT DEFAULT 0,
    source          TEXT,
    discovered_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Approvals
CREATE TABLE IF NOT EXISTS security_approvals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    finding_id      UUID REFERENCES security_findings(id),
    requested_by    UUID,
    status          TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending','approved','rejected')),
    reason          TEXT,
    decided_by      UUID,
    decided_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Audit
CREATE TABLE IF NOT EXISTS security_audit_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type      TEXT NOT NULL,
    actor_id        UUID,
    target_id       TEXT,
    detail          JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index
CREATE INDEX IF NOT EXISTS idx_findings_severity ON security_findings(severity);
CREATE INDEX IF NOT EXISTS idx_findings_status ON security_findings(status);
CREATE INDEX IF NOT EXISTS idx_approvals_status ON security_approvals(status);

