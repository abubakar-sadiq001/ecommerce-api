-- +goose Up
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered','cancelled');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status order_status DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_orders_user_id ON orders (user_id);

CREATE INDEX idx_orders_status ON orders (status);

CREATE INDEX idx_orders_deleted_at ON orders (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS order_status;

DROP INDEX IF EXISTS idx_orders_user_id;

DROP INDEX IF EXISTS idx_orders_status;

DROP INDEX IF EXISTS idx_orders_deleted_at;