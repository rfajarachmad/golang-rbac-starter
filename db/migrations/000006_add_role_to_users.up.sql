-- Add role_id column (nullable first for backfill)
ALTER TABLE users ADD COLUMN role_id INT REFERENCES roles(id);

-- Default all existing users to 'user' role
UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user');

-- Now enforce NOT NULL with a static default (role id 2 = 'user')
ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN role_id SET DEFAULT 2;

-- Seed default admin user (password: admin123)
INSERT INTO users (name, email, password, role_id)
VALUES (
    'Admin',
    'admin@example.com',
    '$2b$12$8zgBT9AC8IxWho98zDxIX.7Dybi1Jx0lFIMd8d.LyxLCXUW4A8oVa',
    (SELECT id FROM roles WHERE name = 'admin')
);
