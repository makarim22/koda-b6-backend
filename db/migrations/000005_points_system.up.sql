ALTER TABLE users ADD COLUMN IF NOT EXISTS points_balance INT DEFAULT 0;

CREATE TABLE IF NOT EXISTS point_ledgers (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL,
    order_id    INT, -- Can be NULL if points were given manually
    points      INT NOT NULL, -- positive for EARN, negative for SPEND
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL
);

ALTER TABLE orders ADD COLUMN IF NOT EXISTS points_used INT DEFAULT 0;
