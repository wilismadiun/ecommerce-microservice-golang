-- +goose Up
-- Perintah untuk membuat tabel stores
CREATE TABLE stores (
    id VARCHAR PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    address TEXT NOT NULL
);

-- +goose Down
-- Perintah untuk membatalkan (rollback) pembuatan tabel stores
DROP TABLE IF EXISTS stores;