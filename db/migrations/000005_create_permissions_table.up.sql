CREATE TABLE IF NOT EXISTS permissions (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

INSERT INTO permissions (name, description) VALUES
    ('user:read', 'View own profile'),
    ('user:update', 'Update own profile'),
    ('user:delete', 'Delete own account / logout'),
    ('contact:create', 'Create contacts'),
    ('contact:read', 'View contacts'),
    ('contact:update', 'Update contacts'),
    ('contact:delete', 'Delete contacts'),
    ('address:create', 'Create addresses'),
    ('address:read', 'View addresses'),
    ('address:update', 'Update addresses'),
    ('address:delete', 'Delete addresses'),
    ('admin:user:list', 'List all users'),
    ('admin:user:read', 'View any user'),
    ('admin:user:update', 'Update any user including role'),
    ('admin:user:delete', 'Delete any user'),
    ('admin:role:manage', 'Create/update/delete roles and permissions');

-- Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'admin';

-- User gets standard CRUD permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'user' AND p.name IN (
    'user:read', 'user:update', 'user:delete',
    'contact:create', 'contact:read', 'contact:update', 'contact:delete',
    'address:create', 'address:read', 'address:update', 'address:delete'
);

-- Viewer gets read-only permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'viewer' AND p.name IN (
    'user:read', 'contact:read', 'address:read'
);
