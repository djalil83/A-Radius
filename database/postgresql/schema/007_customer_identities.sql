-- A-Radius / APB
-- Customer portal identity mapping

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
