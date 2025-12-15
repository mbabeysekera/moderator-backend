-- +goose Up
ALTER TABLE
    products
ADD
    COLUMN added_by BIGINT;

ALTER TABLE
    products
ADD
    CONSTRAINT fk_products_user FOREIGN KEY (added_by) REFERENCES users(id);

-- +goose Down
ALTER TABLE
    products DROP CONSTRAINT IF EXISTS fk_products_user;

ALTER TABLE
    products DROP COLUMN IF EXISTS added_by;