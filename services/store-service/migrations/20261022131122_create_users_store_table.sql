-- +goose Up
-- Perintah untuk membuat tabel user_stores dengan tipe data VARCHAR semua
CREATE TABLE user_stores (
    id VARCHAR(36) PRIMARY KEY,
    store_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(50) NOT NULL
);

-- Menambahkan foreign key constraint ke tabel stores (menggunakan tipe varchar)
ALTER TABLE user_stores 
ADD CONSTRAINT fk_user_stores_store 
FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE;

-- Menambahkan foreign key constraint ke tabel roles (menggunakan tipe varchar)
ALTER TABLE user_stores 
ADD CONSTRAINT fk_user_stores_role 
FOREIGN KEY (role_id) REFERENCES roles(id);

-- +goose Down
-- Perintah untuk membatalkan (rollback) tabel user_stores
DROP TABLE IF EXISTS user_stores;