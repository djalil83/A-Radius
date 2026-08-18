-- A-Radius migration 0002 rollback.
-- Run only after an explicit backup and approval. This removes profile,
-- approval, revision, and audit data created by migration 0002.

BEGIN;

DROP TRIGGER IF EXISTS subscription_profiles_audit_trg ON subscription_profiles;
DROP TRIGGER IF EXISTS subscription_profiles_revision_trg ON subscription_profiles;
DROP TRIGGER IF EXISTS subscription_profiles_updated_at_trg ON subscription_profiles;
DROP TRIGGER IF EXISTS subscription_profile_revisions_no_update_trg ON subscription_profile_revisions;
DROP TRIGGER IF EXISTS audit_events_no_update_trg ON audit_events;
DROP TRIGGER IF EXISTS approval_requests_audit_trg ON approval_requests;

DROP FUNCTION IF EXISTS audit_subscription_profile_change();
DROP FUNCTION IF EXISTS snapshot_subscription_profile();
DROP FUNCTION IF EXISTS set_subscription_profile_updated_at();
DROP FUNCTION IF EXISTS prevent_revision_mutation();
DROP FUNCTION IF EXISTS prevent_audit_mutation();
DROP FUNCTION IF EXISTS audit_approval_request_change();

DROP TABLE IF EXISTS approval_requests;
DROP TABLE IF EXISTS subscription_profile_revisions;
DROP TABLE IF EXISTS audit_events;
DROP TABLE IF EXISTS subscription_profiles;

DROP TYPE IF EXISTS audit_action;
DROP TYPE IF EXISTS approval_request_status;
DROP TYPE IF EXISTS subscription_profile_billing_cycle;
DROP TYPE IF EXISTS subscription_profile_commission_type;
DROP TYPE IF EXISTS subscription_profile_status;
DROP TYPE IF EXISTS subscription_profile_service_type;

COMMIT;
