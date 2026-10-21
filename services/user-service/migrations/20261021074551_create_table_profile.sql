-- +goose Up
CREATE TABLE profiles (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    address TEXT
);

-- +goose Down
DROP TABLE IF EXISTS profiles;