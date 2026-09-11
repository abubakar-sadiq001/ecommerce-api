-- +goose Up
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER UNIQUE NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
    product_id INTEGER UNIQUE NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (cart_id, product_id)
);

CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);

CREATE INDEX idx_cart_items_product_id ON cart_items (product_id);

CREATE INDEX idx_cart_items_deleted_at ON cart_items (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS cart_items;

DROP INDEX IF EXISTS idx_cart_items_cart_id;

DROP INDEX IF EXISTS idx_cart_items_product_id;

DROP INDEX IF EXISTS idx_cart_items_deleted_at;