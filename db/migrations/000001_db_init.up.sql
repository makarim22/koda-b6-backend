
CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    full_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    phone         VARCHAR(20),
    profile_image VARCHAR(255),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS product_category (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS products (
    id                 SERIAL PRIMARY KEY,
    product_name       VARCHAR(150) NOT NULL,
    description        TEXT,
    stock              INT NOT NULL DEFAULT 0,
    base_price         DECIMAL(10,2) NOT NULL,
    category_id        INT REFERENCES product_category(id) ON DELETE SET NULL,
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS variants (
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(100) NOT NULL,
    additional_price DECIMAL(10,2) DEFAULT 0.00,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_variant (
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id INT NOT NULL REFERENCES variants(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, variant_id)
);


CREATE TABLE IF NOT EXISTS sizes (
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(50) NOT NULL,
    additional_price DECIMAL(10,2) DEFAULT 0.00,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_sizes (
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size_id    INT NOT NULL REFERENCES sizes(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, size_id)
);


CREATE TABLE IF NOT EXISTS temperature (
    id         SERIAL PRIMARY KEY,
    label      VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS product_category_map (
    id          SERIAL PRIMARY KEY,
    product_id  INT NOT NULL,
    category_id INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (product_id, category_id),
    FOREIGN KEY (product_id)  REFERENCES products(id)          ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES product_category(id)  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_image (
    id         SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    path       VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orders (
    id           SERIAL PRIMARY KEY,
    customer_id  INT NOT NULL,
    order_date   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    subtotal     DECIMAL(10,2) NOT NULL,
    tax          DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    delivery_fee DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    status       VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE RESTRICT
);


CREATE TABLE IF NOT EXISTS order_detail (
    id             SERIAL PRIMARY KEY,
    order_id       INT NOT NULL,
    product_id     INT NOT NULL,
    size_id        INT,
    temperature_id INT,
    quantity       INT NOT NULL DEFAULT 1,
    unit_price     DECIMAL(10,2) NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id)       REFERENCES orders(id)       ON DELETE CASCADE,
    FOREIGN KEY (product_id)     REFERENCES products(id)     ON DELETE RESTRICT,
    FOREIGN KEY (size_id)        REFERENCES sizes(id)        ON DELETE SET NULL,
    FOREIGN KEY (temperature_id) REFERENCES temperature(id)  ON DELETE SET NULL
);


CREATE TABLE IF NOT EXISTS cart (
    id             SERIAL PRIMARY KEY,
    customer_id    INT NOT NULL,
    product_id     INT NOT NULL,
    size_id        INT,
    temperature_id INT,
    quantity       INT NOT NULL DEFAULT 1,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id)    REFERENCES users(id)        ON DELETE CASCADE,
    FOREIGN KEY (product_id)     REFERENCES products(id)     ON DELETE CASCADE,
    FOREIGN KEY (size_id)        REFERENCES sizes(id)        ON DELETE SET NULL,
    FOREIGN KEY (temperature_id) REFERENCES temperature(id)  ON DELETE SET NULL
);


CREATE TABLE IF NOT EXISTS payment (
    id             SERIAL PRIMARY KEY,
    order_id       INT NOT NULL,
    method         VARCHAR(50) NOT NULL,
    amount         DECIMAL(10,2) NOT NULL,
    status         VARCHAR(50) NOT NULL DEFAULT 'pending',
    transaction_id VARCHAR(150) UNIQUE,
    payment_date   TIMESTAMP,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT
);


CREATE TABLE IF NOT EXISTS user_review (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    product_id INT NOT NULL,
    order_id   INT NOT NULL,
    message    TEXT,
    rating     INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)    REFERENCES users(id)    ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id)   REFERENCES orders(id)   ON DELETE CASCADE
);

