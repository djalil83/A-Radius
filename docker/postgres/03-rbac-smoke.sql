-- Local-only RBAC seed for end-to-end smoke testing.
-- This is not a production identity bootstrap and contains no usable password.

INSERT INTO apb.users (id, username, email, password_hash, full_name, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'api-smoke-test',
    'api-smoke-test@localhost.invalid',
    'disabled-local-smoke-password',
    'API Smoke Test',
    'active'
)
ON CONFLICT (id) DO UPDATE SET status = 'active';

INSERT INTO apb.roles (id, name, description, is_system)
VALUES (
    '00000000-0000-0000-0000-000000000010',
    'api-smoke-admin',
    'Local-only role for profile API smoke testing',
    TRUE
)
ON CONFLICT (id) DO UPDATE SET description = EXCLUDED.description, is_system = TRUE;

INSERT INTO apb.permissions (permission_key, description)
VALUES
    ('subscription_profiles.read', 'Read subscription profiles'),
    ('subscription_profiles.create', 'Create subscription profiles'),
    ('subscription_profiles.update', 'Update subscription profiles'),
    ('subscription_profiles.archive', 'Archive subscription profiles'),
    ('subscription_profiles.read_history', 'Read subscription profile revisions')
ON CONFLICT (permission_key) DO NOTHING;

INSERT INTO apb.user_roles (user_id, role_id)
VALUES ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000010')
ON CONFLICT DO NOTHING;

INSERT INTO apb.role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000010', id
FROM apb.permissions
WHERE permission_key IN (
    'subscription_profiles.read',
    'subscription_profiles.create',
    'subscription_profiles.update',
    'subscription_profiles.archive',
    'subscription_profiles.read_history'
)
ON CONFLICT DO NOTHING;
