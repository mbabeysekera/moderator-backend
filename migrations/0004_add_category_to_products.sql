-- +goose Up
ALTER TABLE
    products
ADD
    COLUMN category TEXT;

CREATE INDEX IF NOT EXISTS idx_product_category ON products(category);

-- +goose Down
ALTER TABLE
    products DROP COLUMN IF EXISTS category