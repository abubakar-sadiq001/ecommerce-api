-- +goose Up
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_carts_user_id ON carts (user_id);

CREATE INDEX idx_carts_deleted_at ON carts (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS carts;

DROP INDEX IF EXISTS idx_carts_user_id;

DROP INDEX IF EXISTS idx_carts_deleted_at;