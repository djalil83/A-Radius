-- A-Radius / APB
-- Customer service foundation

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
