CREATE TABLE IF NOT EXISTS contacts (
    id         VARCHAR(100) PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users(id),
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100),
    email      VARCHAR(100),
    phone      VARCHAR(100),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
