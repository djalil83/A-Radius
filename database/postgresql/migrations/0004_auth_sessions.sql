-- A-Radius / APB
-- Authentication sessions
--
-- Raw session tokens are never stored.
-- Only SHA-256 hashes of session tokens are persisted.


CREATE TABLE IF NOT EXISTS apb.auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,

    token_hash VARCHAR(64) NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    revoked_at TIMESTAMPTZ,

    ip_address INET,
    user_agent TEXT,

    CONSTRAINT fk_auth_sessions_user
        FOREIGN KEY (user_id)
        REFERENCES apb.users(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user
    ON apb.auth_sessions(user_id);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires
    ON apb.auth_sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_active
    ON apb.auth_sessions(user_id, expires_at)
    WHERE revoked_at IS NULL;

INSERT INTO apb.schema_migrations(version)
VALUES ('0004_auth_sessions')
ON CONFLICT (version) DO NOTHING;
