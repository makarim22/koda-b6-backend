CREATE TABLE IF NOT EXISTS product_discount (
                                        "id" SERIAL PRIMARY KEY,
                                        "product_id" INT,
                                        "discount_rate" FLOAT NOT NULL,
                                        "description" VARCHAR(255),
                                        "is_flash_sale" BOOLEAN DEFAULT FALSE,
                                        CONSTRAINT "fk_disc_product" FOREIGN KEY ("product_id") REFERENCES "products"("id") ON DELETE CASCADE
);