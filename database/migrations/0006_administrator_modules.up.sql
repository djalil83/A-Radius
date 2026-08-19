CREATE TABLE IF NOT EXISTS apb.admin_action_proposals (
    id UUID PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    module TEXT NOT NULL CHECK (module IN ('VOUCHER','MITRA','BILLING','PAYMENT','NMS','TEKNISI','FINANCE','SISTEM')),
    action TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    before_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    proposed_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    risk_level TEXT NOT NULL CHECK (risk_level IN ('INFO','LOW','MEDIUM','HIGH','CRITICAL')),
    reason TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'PREVIEW' CHECK (status IN ('PREVIEW','PENDING_APPROVAL','APPROVED','REJECTED','QUEUED','RUNNING','SUCCEEDED','FAILED','ROLLED_BACK')),
    requested_by TEXT NOT NULL,
    approved_by TEXT,
    worker_id TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_admin_action_proposals_branch_status ON apb.admin_action_proposals(branch_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_action_proposals_module ON apb.admin_action_proposals(module, action);

CREATE TABLE IF NOT EXISTS apb.admin_module_audit (
    id BIGSERIAL PRIMARY KEY,
    proposal_id UUID REFERENCES apb.admin_action_proposals(id) ON DELETE SET NULL,
    branch_id BIGINT NOT NULL,
    actor_id TEXT NOT NULL,
    actor_role TEXT NOT NULL,
    module TEXT NOT NULL,
    action TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    result TEXT NOT NULL CHECK (result IN ('PREVIEW','PENDING_APPROVAL','APPROVED','REJECTED','QUEUED','RUNNING','SUCCEEDED','FAILED','ROLLED_BACK')),
    request_id TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_admin_module_audit_branch_time ON apb.admin_module_audit(branch_id, created_at DESC);

CREATE TABLE IF NOT EXISTS apb.admin_ai_reports (
    id UUID PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    module TEXT NOT NULL,
    title TEXT NOT NULL,
    severity TEXT NOT NULL CHECK (severity IN ('INFO','LOW','MEDIUM','HIGH','CRITICAL')),
    finding TEXT NOT NULL,
    recommendation TEXT NOT NULL,
    impact JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN','REVIEWED','DISMISSED','LINKED_TO_PROPOSAL')),
    proposal_id UUID REFERENCES apb.admin_action_proposals(id) ON DELETE SET NULL,
    production_changed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_admin_ai_reports_branch_status ON apb.admin_ai_reports(branch_id, status, created_at DESC);
