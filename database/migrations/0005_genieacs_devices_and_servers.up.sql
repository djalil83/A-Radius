CREATE TABLE IF NOT EXISTS apb.genieacs_servers (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    host TEXT NOT NULL,
    port INTEGER NOT NULL CHECK (port BETWEEN 1 AND 65535),
    username TEXT NOT NULL,
    credential_ref TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE')),
    connection_status TEXT NOT NULL DEFAULT 'UNKNOWN' CHECK (connection_status IN ('CONNECTED','DISCONNECTED','UNKNOWN')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (branch_id, name)
);

CREATE TABLE IF NOT EXISTS apb.onu_devices (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    branch_id BIGINT NOT NULL,
    acs_server_id BIGINT REFERENCES apb.genieacs_servers(id) ON DELETE SET NULL,
    username TEXT,
    pppoe_ip INET,
    tr069_ip INET,
    serial_number TEXT NOT NULL,
    pon_port TEXT,
    manufacturer TEXT,
    model TEXT,
    firmware TEXT,
    ssid TEXT,
    wifi_password_ref TEXT,
    wan_status TEXT NOT NULL DEFAULT 'UNKNOWN',
    last_connected_at TIMESTAMPTZ,
    last_inform_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (branch_id, serial_number)
);

CREATE INDEX IF NOT EXISTS idx_onu_devices_branch_last_inform ON apb.onu_devices(branch_id, last_inform_at);
CREATE INDEX IF NOT EXISTS idx_onu_devices_customer ON apb.onu_devices(customer_id);

CREATE TABLE IF NOT EXISTS apb.genieacs_command_audit (
    id BIGSERIAL PRIMARY KEY,
    actor_id TEXT NOT NULL,
    actor_role TEXT NOT NULL,
    branch_id BIGINT NOT NULL,
    action TEXT NOT NULL,
    module TEXT NOT NULL DEFAULT 'GENIEACS',
    target_type TEXT NOT NULL,
    target_id TEXT NOT NULL,
    request_id TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    result TEXT NOT NULL CHECK (result IN ('PREVIEW','PENDING_APPROVAL','SUCCESS','FAILED','REJECTED')),
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_genieacs_command_audit_branch_time ON apb.genieacs_command_audit(branch_id, created_at DESC);
