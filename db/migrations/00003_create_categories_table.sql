-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_categories_is_active ON categories (is_active);

CREATE INDEX idx_categories_deleted_at ON categories (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS categories;

DROP INDEX IF EXISTS idx_categories_is_active;

DROP INDEX IF EXISTS idx_categories_deleted_at;