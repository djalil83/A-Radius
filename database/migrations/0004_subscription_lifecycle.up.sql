CREATE SCHEMA IF NOT EXISTS apb;

CREATE TABLE IF NOT EXISTS apb.subscription_lifecycle_events (
    id BIGSERIAL PRIMARY KEY,
    subscription_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    source TEXT NOT NULL,
    idempotency_key TEXT NOT NULL UNIQUE,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS apb.subscription_bulk_action_proposals (
    id UUID PRIMARY KEY,
    action TEXT NOT NULL,
    target_filter JSONB NOT NULL DEFAULT '{}'::jsonb,
    target_count INTEGER NOT NULL CHECK (target_count >= 0),
    from_value TEXT,
    to_value TEXT,
    risk TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('PREVIEW','PENDING_APPROVAL','APPROVED','QUEUED','RUNNING','SUCCESS','FAILED','REJECTED')),
    ai_recommendation BOOLEAN NOT NULL DEFAULT false,
    preview_only BOOLEAN NOT NULL DEFAULT true,
    approval_required BOOLEAN NOT NULL DEFAULT true,
    requested_by TEXT NOT NULL,
    approved_by TEXT,
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS apb.subscription_bulk_action_audit (
    id BIGSERIAL PRIMARY KEY,
    proposal_id UUID NOT NULL REFERENCES apb.subscription_bulk_action_proposals(id),
    action TEXT NOT NULL,
    target_count INTEGER NOT NULL,
    requested_by TEXT NOT NULL,
    approved_by TEXT,
    execution TEXT NOT NULL,
    worker_id TEXT,
    error_message TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_subscription_lifecycle_events_subscription ON apb.subscription_lifecycle_events(subscription_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_subscription_bulk_proposals_status ON apb.subscription_bulk_action_proposals(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_subscription_bulk_audit_proposal ON apb.subscription_bulk_action_audit(proposal_id, created_at DESC);
