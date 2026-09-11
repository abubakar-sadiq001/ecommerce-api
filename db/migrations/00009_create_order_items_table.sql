-- +goose Up
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER UNIQUE NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id INTEGER UNIQUE NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE (order_id, product_id)
);

CREATE INDEX idx_order_items_order_id ON order_items (order_id);

CREATE INDEX idx_order_items_product_id ON order_items (product_id);

CREATE INDEX idx_order_items_deleted_at ON order_items (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS order_items;

DROP INDEX IF EXISTS idx_order_items_order_id;

DROP INDEX IF EXISTS idx_order_items_product_id;

DROP INDEX IF EXISTS idx_order_items_deleted_at;