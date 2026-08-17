-- A-Radius / APB
-- Identity foundation
-- Auth/RBAC will build on these tables.

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

CREATE INDEX IF NOT EXISTS idx_users_status
    ON apb.users(status);

CREATE INDEX IF NOT EXISTS idx_users_email
    ON apb.users(email);
