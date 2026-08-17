-- A-Radius / APB
-- Customer foundation

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
