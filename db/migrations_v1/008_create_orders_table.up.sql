CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered','cancelled');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status order_status DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT TIMESTAMP,
        deleted_at TIMESTAMP
    WITH
        TIME ZONE
);

CREATE INDEX idx_order_user_id ON orders (user_id);

CREATE INDEX idx_order_status ON orders (status);

CREATE INDEX idx_orders_deleted_at ON user_id (deleted_at);