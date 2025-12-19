-- +goose Up
ALTER TABLE
    products
ADD
    COLUMN price NUMERIC(10, 2);

-- +goose Down
ALTER TABLE
    products DROP COLUMN IF EXISTS price;