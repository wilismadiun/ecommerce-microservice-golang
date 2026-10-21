-- +goose Up
CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL    
);

-- +goose Down
DROP TABLE IF EXISTS users;