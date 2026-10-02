-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    description TEXT,
    price DECIMAL DEFAULT 0,
    stock INTEGER NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_category_id ON products (category_id);

CREATE INDEX idx_products_sku ON products (sku);

CREATE INDEX idx_products_is_active ON products (is_active);

CREATE INDEX idx_products_deleted_at ON products (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS products;

DROP INDEX IF EXISTS idx_products_category_id;

DROP INDEX IF EXISTS idx_products_sku;

DROP INDEX IF EXISTS idx_products_is_active;

DROP INDEX IF EXISTS idx_products_deleted_at;