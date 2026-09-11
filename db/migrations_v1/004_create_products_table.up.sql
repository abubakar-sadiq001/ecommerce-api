CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    description TEXT,
    price DECIMAL DEFAULT 0,
    sku VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
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

CREATE INDEX idx_products_category_id ON products (category_id);

CREATE INDEX idx_products_sku ON products (sku);

CREATE INDEX idx_products_is_active ON products (is_active);

CREATE INDEX idx_products_deleted_at ON products (deleted_at);