CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    product_id INTEGER UNIQUE NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE,
        UNIQUE (cart_id, product_id)
);

CREATE INDEX idx_order_items_cart_id ON order_items (user_id);

CREATE INDEX idx_order_items_product_id ON order_items (product_id);

CREATE INDEX idx_order_items_deleted_at ON order_items (deleted_at);