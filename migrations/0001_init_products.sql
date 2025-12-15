-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    brand TEXT NOT NULL,
    sku TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    item_code BIGINT NOT NULL,
    in_stock INT NOT NULL,
    image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT uq_product_item UNIQUE(product_id, item_code)
);

CREATE INDEX idx_items_product_id ON items(product_id);

-- +goose Down
ALTER TABLE
    items DROP CONSTRAINT IF EXISTS fk_items_product;

DROP TABLE IF EXISTS items;

DROP TABLE IF EXISTS products;