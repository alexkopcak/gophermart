CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255),
    order_id VARCHAR(255) UNIQUE,
    debet BOOLEAN,
    order_status VARCHAR(255),
    accrual INTEGER,
    uploaded_at TIMESTAMP DEFAULT NOW()
);