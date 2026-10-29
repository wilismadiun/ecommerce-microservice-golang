-- +goose Up
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    store_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    attributes JSONB,
    price DECIMAL(12, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_products_store_id ON products(store_id);

-- +goose Down
DROP INDEX IF EXISTS idx_products_store_id;
DROP TABLE IF EXISTS products;