-- A-Radius migration 0002: subscription profiles, revision history, audit trail.
-- Assumptions:
--   * tenant_id and actor IDs are UUIDs owned by the application identity service.
--   * This migration does not create a users/tenants table because those schemas are project-specific.
--   * External credentials are never stored in these tables.

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE subscription_profile_service_type AS ENUM (
    'FTTH',
    'PPPOE',
    'HOTSPOT_VOUCHER',
    'STATIC_IP'
);

CREATE TYPE subscription_profile_status AS ENUM ('ACTIVE', 'INACTIVE', 'ARCHIVED');
CREATE TYPE subscription_profile_commission_type AS ENUM ('RUPIAH', 'PERCENT');
CREATE TYPE subscription_profile_billing_cycle AS ENUM ('DAILY', 'WEEKLY', 'MONTHLY', 'CUSTOM');
CREATE TYPE approval_request_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED', 'APPLIED');
CREATE TYPE audit_action AS ENUM ('CREATE', 'UPDATE', 'DELETE', 'STATUS_CHANGE', 'EXPORT', 'APPROVAL_REQUESTED', 'APPROVAL_DECIDED', 'APPLIED');

CREATE TABLE subscription_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    service_type subscription_profile_service_type NOT NULL,
    category TEXT,
    media TEXT,
    color CHAR(7) NOT NULL DEFAULT '#1677ff',
    description TEXT,
    status subscription_profile_status NOT NULL DEFAULT 'ACTIVE',

    -- Identifiers/configuration only; never store passwords, tokens, or private keys.
    mikrotik_group TEXT,
    radius_group TEXT,
    rate_limit TEXT,
    upload_bps BIGINT,
    download_bps BIGINT,
    shared_users INTEGER NOT NULL DEFAULT 1,
    vlan_id INTEGER,
    olt_profile TEXT,
    ip_pool TEXT,

    monthly_price BIGINT NOT NULL DEFAULT 0,
    active_days INTEGER NOT NULL DEFAULT 30,
    commission_amount BIGINT NOT NULL DEFAULT 0,
    commission_type subscription_profile_commission_type NOT NULL DEFAULT 'RUPIAH',
    billing_cycle subscription_profile_billing_cycle NOT NULL DEFAULT 'MONTHLY',
    auto_isolate BOOLEAN NOT NULL DEFAULT TRUE,
    billing_note TEXT,

    -- Optimistic locking. The API must include the expected version in UPDATE predicates.
    version BIGINT NOT NULL DEFAULT 1,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT subscription_profiles_color_ck
        CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
    CONSTRAINT subscription_profiles_upload_bps_ck
        CHECK (upload_bps IS NULL OR upload_bps >= 0),
    CONSTRAINT subscription_profiles_download_bps_ck
        CHECK (download_bps IS NULL OR download_bps >= 0),
    CONSTRAINT subscription_profiles_shared_users_ck
        CHECK (shared_users >= 1),
    CONSTRAINT subscription_profiles_vlan_id_ck
        CHECK (vlan_id IS NULL OR vlan_id BETWEEN 1 AND 4094),
    CONSTRAINT subscription_profiles_monthly_price_ck
        CHECK (monthly_price >= 0),
    CONSTRAINT subscription_profiles_active_days_ck
        CHECK (active_days >= 0),
    CONSTRAINT subscription_profiles_commission_amount_ck
        CHECK (commission_amount >= 0),
    CONSTRAINT subscription_profiles_percent_commission_ck
        CHECK (commission_type <> 'PERCENT' OR commission_amount <= 100),
    CONSTRAINT subscription_profiles_version_ck
        CHECK (version >= 1),
    CONSTRAINT subscription_profiles_deleted_status_ck
        CHECK ((deleted_at IS NULL AND status <> 'ARCHIVED') OR (deleted_at IS NOT NULL AND status = 'ARCHIVED'))
);

CREATE UNIQUE INDEX subscription_profiles_tenant_name_uq
    ON subscription_profiles (tenant_id, lower(name))
    WHERE deleted_at IS NULL;

CREATE INDEX subscription_profiles_tenant_status_idx
    ON subscription_profiles (tenant_id, status, updated_at DESC);

CREATE INDEX subscription_profiles_tenant_service_idx
    ON subscription_profiles (tenant_id, service_type, updated_at DESC)
    WHERE deleted_at IS NULL;

CREATE TABLE subscription_profile_revisions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    profile_id UUID NOT NULL REFERENCES subscription_profiles(id) ON DELETE RESTRICT,
    tenant_id UUID NOT NULL,
    version BIGINT NOT NULL,
    operation audit_action NOT NULL,
    snapshot JSONB NOT NULL,
    changed_by UUID,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT subscription_profile_revisions_version_uq UNIQUE (profile_id, version)
);

CREATE INDEX subscription_profile_revisions_profile_idx
    ON subscription_profile_revisions (profile_id, version DESC);

CREATE TABLE audit_events (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id UUID NOT NULL,
    actor_id UUID,
    request_id UUID,
    action audit_action NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id UUID,
    entity_version BIGINT,
    before_data JSONB,
    after_data JSONB,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT audit_events_entity_type_ck
        CHECK (length(trim(entity_type)) BETWEEN 1 AND 100),
    CONSTRAINT audit_events_payload_ck
        CHECK (before_data IS NOT NULL OR after_data IS NOT NULL OR metadata <> '{}'::jsonb)
);

CREATE INDEX audit_events_tenant_created_idx
    ON audit_events (tenant_id, created_at DESC);

CREATE INDEX audit_events_entity_idx
    ON audit_events (entity_type, entity_id, created_at DESC);

CREATE INDEX audit_events_request_idx
    ON audit_events (request_id)
    WHERE request_id IS NOT NULL;

CREATE TABLE approval_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    profile_id UUID NOT NULL REFERENCES subscription_profiles(id) ON DELETE RESTRICT,
    requested_by UUID NOT NULL,
    decided_by UUID,
    status approval_request_status NOT NULL DEFAULT 'PENDING',
    requested_version BIGINT NOT NULL,
    requested_changes JSONB NOT NULL,
    decision_note TEXT,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    decided_at TIMESTAMPTZ,
    applied_at TIMESTAMPTZ,
    CONSTRAINT approval_requests_requested_version_ck CHECK (requested_version >= 1),
    CONSTRAINT approval_requests_decision_consistency_ck CHECK (
        (status = 'PENDING' AND decided_by IS NULL AND decided_at IS NULL AND applied_at IS NULL)
        OR (status IN ('APPROVED', 'REJECTED', 'CANCELLED') AND decided_by IS NOT NULL AND decided_at IS NOT NULL AND applied_at IS NULL)
        OR (status = 'APPLIED' AND decided_by IS NOT NULL AND decided_at IS NOT NULL AND applied_at IS NOT NULL)
    )
);

CREATE UNIQUE INDEX approval_requests_one_pending_profile_uq
    ON approval_requests (profile_id)
    WHERE status = 'PENDING';

CREATE INDEX approval_requests_tenant_status_idx
    ON approval_requests (tenant_id, status, requested_at DESC);

CREATE INDEX approval_requests_profile_idx
    ON approval_requests (profile_id, requested_at DESC);

CREATE OR REPLACE FUNCTION set_subscription_profile_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    IF NEW.version <> OLD.version + 1 THEN
        RAISE EXCEPTION 'subscription profile version must increment by exactly one (old %, new %)', OLD.version, NEW.version
            USING ERRCODE = '22000';
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER subscription_profiles_updated_at_trg
BEFORE UPDATE ON subscription_profiles
FOR EACH ROW EXECUTE FUNCTION set_subscription_profile_updated_at();

CREATE OR REPLACE FUNCTION prevent_revision_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'subscription_profile_revisions is append-only';
END;
$$;

CREATE TRIGGER subscription_profile_revisions_no_update_trg
BEFORE UPDATE OR DELETE ON subscription_profile_revisions
FOR EACH ROW EXECUTE FUNCTION prevent_revision_mutation();

CREATE OR REPLACE FUNCTION snapshot_subscription_profile()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO subscription_profile_revisions
            (profile_id, tenant_id, version, operation, snapshot, changed_by)
        VALUES
            (NEW.id, NEW.tenant_id, NEW.version, 'CREATE', to_jsonb(NEW), NEW.created_by);
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO subscription_profile_revisions
            (profile_id, tenant_id, version, operation, snapshot, changed_by)
        VALUES
            (NEW.id, NEW.tenant_id, NEW.version,
             CASE WHEN OLD.status IS DISTINCT FROM NEW.status THEN 'STATUS_CHANGE' ELSE 'UPDATE' END,
             to_jsonb(NEW), NEW.updated_by);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$;

CREATE TRIGGER subscription_profiles_revision_trg
after INSERT OR UPDATE ON subscription_profiles
FOR EACH ROW EXECUTE FUNCTION snapshot_subscription_profile();

CREATE OR REPLACE FUNCTION prevent_audit_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'audit_events is append-only';
END;
$$;

CREATE TRIGGER audit_events_no_update_trg
BEFORE UPDATE OR DELETE ON audit_events
FOR EACH ROW EXECUTE FUNCTION prevent_audit_mutation();

CREATE OR REPLACE FUNCTION audit_subscription_profile_change()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    action_value audit_action;
    actor_value UUID;
    request_value UUID;
BEGIN
    request_value := NULLIF(current_setting('app.request_id', true), '')::uuid;
    IF TG_OP = 'INSERT' THEN
        action_value := 'CREATE';
        actor_value := NEW.created_by;
        INSERT INTO audit_events (tenant_id, actor_id, request_id, action, entity_type, entity_id, entity_version, after_data)
        VALUES (NEW.tenant_id, actor_value, request_value, action_value, 'subscription_profile', NEW.id, NEW.version, to_jsonb(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        action_value := CASE WHEN OLD.status IS DISTINCT FROM NEW.status THEN 'STATUS_CHANGE' ELSE 'UPDATE' END;
        actor_value := NEW.updated_by;
        INSERT INTO audit_events (tenant_id, actor_id, request_id, action, entity_type, entity_id, entity_version, before_data, after_data)
        VALUES (NEW.tenant_id, actor_value, request_value, action_value, 'subscription_profile', NEW.id, NEW.version, to_jsonb(OLD), to_jsonb(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit_events (tenant_id, actor_id, request_id, action, entity_type, entity_id, entity_version, before_data)
        VALUES (OLD.tenant_id, NULLIF(current_setting('app.actor_id', true), '')::uuid, request_value, 'DELETE', 'subscription_profile', OLD.id, OLD.version, to_jsonb(OLD));
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$;

CREATE TRIGGER subscription_profiles_audit_trg
after INSERT OR UPDATE OR DELETE ON subscription_profiles
FOR EACH ROW EXECUTE FUNCTION audit_subscription_profile_change();

CREATE OR REPLACE FUNCTION audit_approval_request_change()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    action_value audit_action;
    actor_value UUID;
BEGIN
    IF TG_OP = 'INSERT' THEN
        action_value := 'APPROVAL_REQUESTED';
        actor_value := NEW.requested_by;
        INSERT INTO audit_events (tenant_id, actor_id, action, entity_type, entity_id, entity_version, after_data)
        VALUES (NEW.tenant_id, actor_value, action_value, 'approval_request', NEW.id, NEW.requested_version, to_jsonb(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        action_value := CASE WHEN NEW.status IS DISTINCT FROM OLD.status THEN 'APPROVAL_DECIDED' ELSE 'UPDATE' END;
        actor_value := COALESCE(NEW.decided_by, NEW.requested_by);
        INSERT INTO audit_events (tenant_id, actor_id, action, entity_type, entity_id, entity_version, before_data, after_data)
        VALUES (NEW.tenant_id, actor_value, action_value, 'approval_request', NEW.id, NEW.requested_version, to_jsonb(OLD), to_jsonb(NEW));
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$;

CREATE TRIGGER approval_requests_audit_trg
after INSERT OR UPDATE ON approval_requests
FOR EACH ROW EXECUTE FUNCTION audit_approval_request_change();

COMMENT ON TABLE audit_events IS 'Append-only security and business audit trail. Do not store secrets in JSON payloads.';
COMMENT ON TABLE subscription_profile_revisions IS 'Immutable snapshots for profile version history and rollback review.';
COMMENT ON COLUMN subscription_profiles.version IS 'Optimistic-lock version; API updates must match the expected version.';
COMMENT ON COLUMN subscription_profiles.rate_limit IS 'MikroTik/RADIUS policy identifier or value; never a credential.';

COMMIT;
