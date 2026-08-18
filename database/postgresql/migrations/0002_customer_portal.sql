-- A-Radius / APB
-- Customer portal RBAC permission

BEGIN;

INSERT INTO apb.permissions (
    permission_key,
    description
)
VALUES (
    'customer.portal.read',
    'Allow authenticated customers to read their own customer portal'
)
ON CONFLICT (permission_key) DO NOTHING;

INSERT INTO apb.roles (
    name,
    description,
    is_system
)
VALUES (
    'customer',
    'Customer portal role',
    TRUE
)
ON CONFLICT (name) DO NOTHING;

INSERT INTO apb.role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM apb.roles r
JOIN apb.permissions p
    ON p.permission_key = 'customer.portal.read'
WHERE r.name = 'customer'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO apb.schema_migrations(version)
VALUES ('0002_customer_portal')
ON CONFLICT (version) DO NOTHING;

COMMIT;
