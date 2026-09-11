CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER UNIQUE NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
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

CREATE INDEX idx_cart_items_cart_id ON cart_items (user_id);

CREATE INDEX idx_cart_items_product_id ON cart_items (product_id);

CREATE INDEX idx_cart_items_deleted_at ON cart_items (deleted_at);