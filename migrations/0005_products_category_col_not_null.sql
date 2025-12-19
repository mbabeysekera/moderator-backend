-- +goose Up
ALTER TABLE
    products
ALTER COLUMN
    category
SET
    NOT NULL;

-- +goose Down
ALTER TABLE
    products
ALTER COLUMN
    category TEXT NULL;