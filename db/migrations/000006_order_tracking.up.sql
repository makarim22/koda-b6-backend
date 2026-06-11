CREATE TABLE IF NOT EXISTS order_tracking (
    id          SERIAL PRIMARY KEY,
    order_id    INT NOT NULL,
    status      VARCHAR(50) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_order_tracking_order_id ON order_tracking(order_id);
