CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);

INSERT INTO roles (name, description, permissions) VALUES
    ('super_admin', 'Full access', '["*"]'),
    ('cinema_admin', 'Manage cinema operations', '["cinema:*","movie:*","showtime:*"]'),
    ('cashier', 'Validate tickets and handle walk-in', '["ticket:validate","booking:read"]'),
    ('customer', 'Purchase tickets', '["booking:create","booking:read","payment:create"]');