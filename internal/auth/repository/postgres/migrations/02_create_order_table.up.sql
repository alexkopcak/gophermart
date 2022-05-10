CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    order_id VARCHAR(255) UNIQUE,
    debet BOOLEAN,
    order_status VARCHAR(255),
    accrual INTEGER,
    uploaded_at TIMESTAMP DEFAULT NOW()
);