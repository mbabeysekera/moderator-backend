-- +goose Up
-- 1. Create tenants table
CREATE TABLE tenants (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- 2. Insert default tenant (Migration of existing data)
INSERT INTO tenants (id, name, slug) VALUES (1, 'VA Collection', 'va-collection');
-- Ensure the sequence is updated to avoid conflicts
SELECT setval(pg_get_serial_sequence('tenants', 'id'), 1);

-- 3. Create tenant_users table
CREATE TABLE tenant_users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role TEXT DEFAULT 'USER' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (tenant_id, user_id),
    CONSTRAINT fk_tenant_users_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_tenant_users_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

-- 4. Migrate existing users to default tenant
INSERT INTO tenant_users (tenant_id, user_id, role)
SELECT 1, id, role FROM users;

-- 5. Add tenant_id to products and migrate data
ALTER TABLE products
ADD COLUMN tenant_id BIGINT DEFAULT 1; -- Set default for new rows temporarily helps with not null, but we will update.

-- Update existing rows
UPDATE products SET tenant_id = 1 WHERE tenant_id IS NULL;

-- Enforce Foreign Key
ALTER TABLE products
ADD CONSTRAINT fk_products_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON DELETE CASCADE;

ALTER TABLE products ALTER COLUMN tenant_id SET NOT NULL;

CREATE INDEX idx_products_tenant_id ON products(tenant_id);

-- 6. Add tenant_id to items and migrate data
ALTER TABLE items
ADD COLUMN tenant_id BIGINT DEFAULT 1;

-- Update existing rows
UPDATE items SET tenant_id = 1 WHERE tenant_id IS NULL;

-- Enforce Foreign Key
ALTER TABLE items
ADD CONSTRAINT fk_items_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON DELETE CASCADE;

ALTER TABLE items ALTER COLUMN tenant_id SET NOT NULL;

CREATE INDEX idx_items_tenant_id ON items(tenant_id);


-- +goose Down
-- Reverting changes
DROP INDEX IF EXISTS idx_items_tenant_id;
ALTER TABLE items DROP CONSTRAINT IF EXISTS fk_items_tenant;
ALTER TABLE items DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_products_tenant_id;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_tenant;
ALTER TABLE products DROP COLUMN IF EXISTS tenant_id;

DROP TABLE IF EXISTS tenant_users;
DROP TABLE IF EXISTS tenants;
