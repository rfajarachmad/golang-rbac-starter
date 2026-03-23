CREATE TABLE IF NOT EXISTS addresses (
    id          VARCHAR(100) PRIMARY KEY,
    contact_id  VARCHAR(100) NOT NULL REFERENCES contacts(id),
    street      VARCHAR(255),
    city        VARCHAR(100),
    province    VARCHAR(100),
    postal_code VARCHAR(20),
    country     VARCHAR(100),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
