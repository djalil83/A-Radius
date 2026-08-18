-- A-Radius / APB
-- Customer core migration
--
-- Adds customer, service, and customer identity tables
-- required by the Customer Portal.

BEGIN;

CREATE TABLE IF NOT EXISTS apb.customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(30),
    address TEXT,
    village VARCHAR(100),
    district VARCHAR(100),
    regency VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    latitude NUMERIC(10,7),
    longitude NUMERIC(10,7),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT customers_status_check
        CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE INDEX IF NOT EXISTS idx_customers_name
    ON apb.customers(name);

CREATE INDEX IF NOT EXISTS idx_customers_status
    ON apb.customers(status);

CREATE INDEX IF NOT EXISTS idx_customers_location
    ON apb.customers(latitude, longitude);


CREATE TABLE IF NOT EXISTS apb.services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL,
    service_code VARCHAR(80) NOT NULL UNIQUE,
    service_type VARCHAR(30) NOT NULL,
    username VARCHAR(150),
    package_name VARCHAR(150),
    download_speed BIGINT,
    upload_speed BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    activated_at TIMESTAMPTZ,
    suspended_at TIMESTAMPTZ,
    terminated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_services_customer
        FOREIGN KEY (customer_id)
        REFERENCES apb.customers(id)
        ON DELETE RESTRICT,

    CONSTRAINT services_type_check
        CHECK (
            service_type IN (
                'pppoe',
                'hotspot',
                'static',
                'ftth',
                'wireless',
                'other'
            )
        ),

    CONSTRAINT services_status_check
        CHECK (
            status IN (
                'active',
                'inactive',
                'suspended',
                'terminated'
            )
        )
);

CREATE INDEX IF NOT EXISTS idx_services_customer
    ON apb.services(customer_id);

CREATE INDEX IF NOT EXISTS idx_services_type
    ON apb.services(service_type);

CREATE INDEX IF NOT EXISTS idx_services_status
    ON apb.services(status);


CREATE TABLE IF NOT EXISTS apb.customer_identities (
    user_id UUID PRIMARY KEY,
    customer_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_customer_identity_user
        FOREIGN KEY (user_id)
        REFERENCES apb.users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_customer_identity_customer
        FOREIGN KEY (customer_id)
        REFERENCES apb.customers(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_customer_identities_customer
    ON apb.customer_identities(customer_id);


INSERT INTO apb.schema_migrations(version)
VALUES ('0003_customer_core')
ON CONFLICT (version) DO NOTHING;

COMMIT;
