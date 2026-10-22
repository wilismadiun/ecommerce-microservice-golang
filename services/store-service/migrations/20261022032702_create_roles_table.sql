-- +goose Up
CREATE TABLE roles (
    id VARCHAR PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- Masukkan data master langsung di sini
INSERT INTO roles (id, name) VALUES (1, 'OWNER');
INSERT INTO roles (id, name) VALUES (2, 'ADMIN');

-- +goose Down
DROP TABLE roles;