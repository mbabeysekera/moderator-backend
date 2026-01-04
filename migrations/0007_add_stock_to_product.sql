-- +goose Up
ALTER TABLE
    products
ADD
    COLUMN in_stock INT NOT NULL DEFAULT 0;

UPDATE
    products p
SET
    in_stock = sub.avg_stock
FROM
    (
        SELECT
            product_id,
            CASE
                WHEN COUNT(*) = 0 THEN 0
                ELSE (SUM(in_stock) / COUNT(*))
            END AS avg_stock
        FROM
            items
        GROUP BY
            product_id
    ) sub
WHERE
    p.id = sub.product_id;

ALTER TABLE
    items DROP CONSTRAINT IF EXISTS uq_product_item;

ALTER TABLE
    items DROP COLUMN in_stock;

ALTER TABLE
    items DROP COLUMN item_code;

-- +goose Down
ALTER TABLE
    items
ADD
    COLUMN item_code BIGINT NOT NULL DEFAULT 1;

ALTER TABLE
    items
ADD
    COLUMN in_stock INT NOT NULL DEFAULT 0;

UPDATE
    items i
SET
    in_stock = p.in_stock
FROM
    products p
WHERE
    i.product_id = p.id;

ALTER TABLE
    items
ADD
    CONSTRAINT uq_product_item UNIQUE (product_id, item_code);

ALTER TABLE
    products DROP COLUMN in_stock;