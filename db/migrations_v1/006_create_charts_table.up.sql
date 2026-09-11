CREATE TABLE charts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
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

CREATE INDEX idx_charts_user_id ON products_images (user_id);

CREATE INDEX idx_charts_deleted_at ON products_images (deleted_at);