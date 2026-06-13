-- Insert categories
TRUNCATE TABLE product_category_map CASCADE;
-- product_category has foreign keys to product_category_map, we cascaded it.
-- but since product_category is referenced by product_category_map, we might want to TRUNCATE product_category CASCADE as well
TRUNCATE TABLE product_category CASCADE;

INSERT INTO product_category (id, name, description) VALUES
(1, 'Favorite Product', 'Top picked products'),
(2, 'Coffee', 'Coffee beverages'),
(3, 'Non Coffee', 'Non coffee beverages'),
(4, 'Foods', 'Food and snacks'),
(5, 'Add-On', 'Extra items'),
(6, 'Flash Sale', 'Products on flash sale'),
(7, 'Buy 1 Get 1', 'Buy 1 get 1 free promotions'),
(8, 'Birthday Package', 'Birthday special packages');

-- Reset sequence if needed
SELECT setval('product_category_id_seq', (SELECT MAX(id) FROM product_category));

-- Map Products 1-20 to 'Coffee' (id: 2)
INSERT INTO product_category_map (product_id, category_id)
SELECT id, 2 FROM products WHERE id BETWEEN 1 AND 20;

-- Map Products 21-35 to 'Foods' (id: 4)
INSERT INTO product_category_map (product_id, category_id)
SELECT id, 4 FROM products WHERE id BETWEEN 21 AND 35;

-- Randomly map some products to 'Favorite Product' (id: 1)
INSERT INTO product_category_map (product_id, category_id) VALUES
(1, 1), (3, 1), (5, 1), (12, 1), (21, 1), (30, 1), (35, 1);

-- Randomly map some products to 'Flash Sale' (id: 6)
INSERT INTO product_category_map (product_id, category_id) VALUES
(2, 6), (7, 6), (15, 6), (22, 6), (25, 6);

-- Randomly map some products to 'Buy 1 Get 1' (id: 7)
INSERT INTO product_category_map (product_id, category_id) VALUES
(4, 7), (10, 7), (24, 7), (31, 7);

-- Randomly map some products to 'Birthday Package' (id: 8)
INSERT INTO product_category_map (product_id, category_id) VALUES
(6, 8), (14, 8), (28, 8), (33, 8);
