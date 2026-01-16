-- +goose Up
-- 1. Create app table
CREATE TABLE app (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

-- 2. Insert default app (Migration of existing data)
INSERT INTO app (id, name, slug) VALUES (1, 'VA Collection', 'va-collection');
-- Ensure the sequence is updated to avoid conflicts
SELECT setval(pg_get_serial_sequence('app', 'id'), 1);

-- 3. Update Users Table
ALTER TABLE users
ADD COLUMN app_id BIGINT DEFAULT 1;

ALTER TABLE users
ADD CONSTRAINT fk_users_app
FOREIGN KEY (app_id) REFERENCES app(id)
ON DELETE CASCADE;

ALTER TABLE users ALTER COLUMN app_id SET NOT NULL;
CREATE INDEX idx_users_app_id ON users(app_id);

-- 4. Update Products Table
ALTER TABLE products
ADD COLUMN app_id BIGINT DEFAULT 1;

ALTER TABLE products
ADD CONSTRAINT fk_products_app
FOREIGN KEY (app_id) REFERENCES app(id)
ON DELETE CASCADE;

ALTER TABLE products ALTER COLUMN app_id SET NOT NULL;
CREATE INDEX idx_products_app_id ON products(app_id);

-- 5. Drop old global SKU constraint and add app-scoped unique constraint
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_sku_key;

CREATE UNIQUE INDEX uq_products_app_sku
ON products (app_id, sku);

-- 6. Update Items Table
ALTER TABLE items
ADD COLUMN app_id BIGINT DEFAULT 1;

ALTER TABLE items
ADD CONSTRAINT fk_items_app
FOREIGN KEY (app_id) REFERENCES app(id)
ON DELETE CASCADE;

ALTER TABLE items ALTER COLUMN app_id SET NOT NULL;
CREATE INDEX idx_items_app_id ON items(app_id);


-- +goose Down
-- Revert Items changes
DROP INDEX IF EXISTS idx_items_app_id;
ALTER TABLE items DROP CONSTRAINT IF EXISTS fk_items_app;
ALTER TABLE items DROP COLUMN IF EXISTS app_id;

-- Revert Products changes
DROP INDEX IF EXISTS uq_products_app_sku;
-- Attempt to restore the original unique constraint (might fail if data is duplicate across apps now, but best effort for rollback)
-- In a real rollback scenario, duplicates would need to be handled manually.
ALTER TABLE products ADD CONSTRAINT products_sku_key UNIQUE (sku);

DROP INDEX IF EXISTS idx_products_app_id;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_app;
ALTER TABLE products DROP COLUMN IF EXISTS app_id;

-- Revert Users changes
DROP INDEX IF EXISTS idx_users_app_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_app;
ALTER TABLE users DROP COLUMN IF EXISTS app_id;

-- Drop App table
DROP TABLE IF EXISTS app;
