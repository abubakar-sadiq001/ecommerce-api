-- +goose Up
CREATE TABLE products_images (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    product_id INTEGER NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_images_product_id ON products_images (product_id);

CREATE INDEX idx_products_is_primary ON products_images (is_primary);

-- CREATE INDEX idx_products_deleted_at ON products_images (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS products_images;

DROP INDEX IF EXISTS idx_products_images_product_id;

DROP INDEX IF EXISTS idx_products_is_primary;

DROP INDEX IF EXISTS idx_products_deleted_at;