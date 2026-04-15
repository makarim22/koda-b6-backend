package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/cache"


	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
)

type ProductRepository struct {
	db *pgxpool.Pool
	cache *cache.RedisCache
}

func NewProductRepository(db *pgxpool.Pool, cache *cache.RedisCache) *ProductRepository {
	return &ProductRepository{
		db: db,
		cache : cache,
	}
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	cacheKey := "products:all"

	var cachedProducts []models.Product
	if err := p.cache.Get(ctx, cacheKey, &cachedProducts); err == nil {
		log.Printf("✅ [CACHE HIT] Retrieved %d products from cache", len(cachedProducts))
		return cachedProducts, nil  // Cache HIT
	}
	log.Printf("❌ [CACHE MISS] Cache key '%s' not found, querying database", cacheKey)

	// Query 1: Get all products
	log.Printf("📡 [QUERY 1] Starting to fetch all products...")
	productRows, err := p.db.Query(ctx, `
		SELECT 
			id,
			product_name,
			description,
			stock,
			base_price
		FROM products
		ORDER BY id DESC
	`)
	if err != nil {
		log.Printf("❌ [QUERY 1 ERROR] Failed to query products: %v", err)
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer productRows.Close()

	products := []models.Product{}
	productIDs := []int{}

	for productRows.Next() {
		var product models.Product
		err := productRows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Description,
			&product.Stock,
			&product.BasePrice,
		)
		if err != nil {
			log.Printf("❌ [SCAN ERROR] Failed to scan product: %v", err)
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		product.Images = []models.ProductImage{}
		product.Variants = []models.Variant{}
		product.Sizes = []models.Size{}

		products = append(products, product)
		productIDs = append(productIDs, product.ID)
	}

	if err = productRows.Err(); err != nil {
		log.Printf("❌ [ITERATION ERROR] Error iterating products: %v", err)
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	log.Printf("✅ [QUERY 1 COMPLETE] Found %d products with IDs: %v", len(products), productIDs)

	// Handle empty products
	if len(products) == 0 {
		log.Printf("⚠️ [EMPTY RESULT] No products found in database")
		// Cache empty result (5 menit)
		if err := p.cache.Set(ctx, cacheKey, []models.Product{}, 5*time.Minute); err != nil {
			log.Printf("⚠️ [CACHE SET ERROR] Failed to cache empty result: %v", err)
		}
		return products, nil
	}

	// Query 2: Get all images for these products
	log.Printf("📡 [QUERY 2] Starting to fetch images for %d products...", len(products))
	imageRows, err := p.db.Query(ctx, `
		SELECT 
			id,
			product_id,
			path,
			is_primary
		FROM product_image
		WHERE product_id = ANY($1)
		ORDER BY product_id, is_primary DESC
	`, productIDs)
	if err != nil {
		log.Printf("❌ [QUERY 2 ERROR] Failed to query product images: %v", err)
		return nil, fmt.Errorf("failed to query product images: %w", err)
	}
	defer imageRows.Close()

	// Create a map for quick product lookup
	productMap := make(map[int]*models.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	imageCount := 0
	for imageRows.Next() {
		var (
			id        int
			productID int
			path      string
			isPrimary bool
		)
		err := imageRows.Scan(&id, &productID, &path, &isPrimary)
		if err != nil {
			log.Printf("❌ [SCAN IMAGE ERROR] Failed to scan image: %v", err)
			return nil, fmt.Errorf("failed to scan image: %w", err)
		}

		if product, exists := productMap[productID]; exists {
			product.Images = append(product.Images, models.ProductImage{
				ID:        id,
				ProductID: productID,
				Path:      path,
				IsPrimary: isPrimary,
			})
			imageCount++
		}
	}

	if err = imageRows.Err(); err != nil {
		log.Printf("❌ [ITERATION IMAGE ERROR] Error iterating images: %v", err)
		return nil, fmt.Errorf("error iterating images: %w", err)
	}

	log.Printf("✅ [QUERY 2 COMPLETE] Loaded %d images for products", imageCount)

	// Query 3: Get all variants for these products
	log.Printf("📡 [QUERY 3] Starting to fetch variants for %d products...", len(products))
	variantRows, err := p.db.Query(ctx, `
		SELECT DISTINCT
			pv.variant_id,
			v.name,
			pv.product_id
		FROM product_variant pv
		JOIN variants v ON pv.variant_id = v.id
		WHERE pv.product_id = ANY($1)
		ORDER BY pv.product_id
	`, productIDs)
	if err != nil {
		log.Printf("❌ [QUERY 3 ERROR] Failed to query variants: %v", err)
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer variantRows.Close()

	variantCount := 0
	for variantRows.Next() {
		var (
			variantID int
			name      string
			productID int
		)
		err := variantRows.Scan(&variantID, &name, &productID)
		if err != nil {
			log.Printf("❌ [SCAN VARIANT ERROR] Failed to scan variant: %v", err)
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		if product, exists := productMap[productID]; exists {
			product.Variants = append(product.Variants, models.Variant{
				ID:   variantID,
				Name: name,
			})
			variantCount++
		}
	}

	if err = variantRows.Err(); err != nil {
		log.Printf("❌ [ITERATION VARIANT ERROR] Error iterating variants: %v", err)
		return nil, fmt.Errorf("error iterating variants: %w", err)
	}

	log.Printf("✅ [QUERY 3 COMPLETE] Loaded %d variants for products", variantCount)

	// Query 4: Get all sizes for these products
	log.Printf("📡 [QUERY 4] Starting to fetch sizes for %d products...", len(products))
	sizeRows, err := p.db.Query(ctx, `
		SELECT DISTINCT
			ps.size_id,
			s.name,
			ps.product_id
		FROM product_sizes ps
		JOIN sizes s ON ps.size_id = s.id
		WHERE ps.product_id = ANY($1)
		ORDER BY ps.product_id
	`, productIDs)
	if err != nil {
		log.Printf("❌ [QUERY 4 ERROR] Failed to query sizes: %v", err)
		return nil, fmt.Errorf("failed to query sizes: %w", err)
	}
	defer sizeRows.Close()

	sizeCount := 0
	for sizeRows.Next() {
		var (
			sizeID    int
			name      string
			productID int
		)
		err := sizeRows.Scan(&sizeID, &name, &productID)
		if err != nil {
			log.Printf("❌ [SCAN SIZE ERROR] Failed to scan size: %v", err)
			return nil, fmt.Errorf("failed to scan size: %w", err)
		}

		if product, exists := productMap[productID]; exists {
			product.Sizes = append(product.Sizes, models.Size{
				ID:   sizeID,
				Name: name,
			})
			sizeCount++
		}
	}

	if err = sizeRows.Err(); err != nil {
		log.Printf("❌ [ITERATION SIZE ERROR] Error iterating sizes: %v", err)
		return nil, fmt.Errorf("error iterating sizes: %w", err)
	}

	log.Printf("✅ [QUERY 4 COMPLETE] Loaded %d sizes for products", sizeCount)

	// Cache the results
	log.Printf("💾 [CACHE SET] Caching %d products with key '%s' for 10 minutes", len(products), cacheKey)
	if err := p.cache.Set(ctx, cacheKey, products, 10*time.Minute); err != nil {
		log.Printf("⚠️ [CACHE SET ERROR] Failed to cache products: %v", err)  // Log saja, jangan return error
	}

	log.Printf("✅ [FINAL RESULT] Successfully returning %d products", len(products))
	return products, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product

	err := p.db.QueryRow(ctx,
		`SELECT id, product_name, description, stock, base_price FROM products WHERE id = $1`,
		id).Scan(&product.ID, &product.ProductName, &product.Description, &product.Stock, &product.BasePrice)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return &product, nil
}

func (p *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	// Start a transaction
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// query 1
	err = tx.QueryRow(ctx,
		`INSERT INTO products (product_name, description, stock, base_price)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		product.ProductName, product.Description, product.Stock, product.BasePrice).
		Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	// query 2
	if len(product.Variants) > 0 {
		for _, variant := range product.Variants {
			_, err := tx.Exec(ctx,
				`INSERT INTO product_variant (product_id, variant_id)
				 VALUES ($1, $2)`,
				product.ID, variant.ID)
			
			if err != nil {
				return fmt.Errorf("failed to create product variant: %w", err)
			}
		}
	}

	// query 3
	if len(product.Sizes) > 0 {
		for _, size := range product.Sizes {
			_, err := tx.Exec(ctx,
				`INSERT INTO product_sizes (product_id, size_id)
				 VALUES ($1, $2)`,
				product.ID, size.ID)
			
			if err != nil {
				return fmt.Errorf("failed to create product size: %w", err)
			}
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	err := p.db.QueryRow(ctx,
		`UPDATE products 
		 SET product_name = $1, description = $2, stock = $3, base_price = $4
		 WHERE id = $5
		 RETURNING id, product_name, description, stock, base_price`,
		product.ProductName, product.Description, product.Stock, product.BasePrice, product.ID).
		Scan(&product.ID, &product.ProductName, &product.Description, &product.Stock, &product.BasePrice)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("product with ID %d not found", product.ID)
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (p *ProductRepository) Delete(ctx context.Context, id int) error {
	commandTag, err := p.db.Exec(ctx,
		`DELETE FROM products WHERE id = $1`,
		id)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

func (p *ProductRepository) MostReview(ctx context.Context) (*[]models.Product, error) {
	var products []models.Product

	query :=
		`SELECT p.id, p.product_name, p.description, p.base_price, count(ur.product_id) 
         FROM products p 
         join user_review ur 
         on p.id = ur.product_id
         group by p.id, p.product_name, p.description, p.base_price
         order by count(ur.product_id) desc`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying most reviewed products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		var reviewCount int
		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Description,
			&product.BasePrice,
			&reviewCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through product rows: %w", err)
	}

	return &products, nil
}

func (p *ProductRepository) UpdateStock(ctx context.Context, id, quantity int) error {
	query := `UPDATE products SET stock = stock - $1 WHERE id = $2`
	_, err := p.db.Exec(ctx, query, quantity, id)
	return err
}

func (p *ProductRepository) MostReviewWithPrimaryImage(ctx context.Context) (*[]models.ProductWithImages, error) {
	var products []models.ProductWithImages

	query := `
	SELECT 
		p.id, 
		p.product_name, 
		p.description, 
		p.base_price,
		pi.id as image_id,
		pi.path,
		pi.is_primary
	FROM products p 
	LEFT JOIN user_review ur ON p.id = ur.product_id
	LEFT JOIN product_image pi ON p.id = pi.product_id AND pi.is_primary = true
	GROUP BY p.id, p.product_name, p.description, p.base_price, pi.id, pi.path, pi.is_primary
	ORDER BY COUNT(ur.product_id) DESC
	`

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying most reviewed products with primary image: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductWithImages
		var imageID *int
		var imagePath *string
		var isPrimary *bool

		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Description,
			&product.BasePrice,
			&imageID,
			&imagePath,
			&isPrimary,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}

		// Set image jika ada
		if imageID != nil && imagePath != nil {
			product.Images = []models.ProductImage{
				{
					ID:        *imageID,
					ProductID: product.ID,
					Path:      *imagePath,
					IsPrimary: *isPrimary,
				},
			}
		} else {
			product.Images = []models.ProductImage{}
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through product rows: %w", err)
	}

	return &products, nil
}

func (r *OrderRepository) GetOrderDetails(ctx context.Context, orderID int) ([]models.OrderDetailResponse, error) {
	log.Printf("[GetOrderDetails] Starting query for orderID: %d", orderID)

	query := `
		SELECT 
			oi.id,
			oi.product_id,
			p.product_name,
			oi.quantity,
			p.base_price,
			oi.size_id,
			s.name as size_name,
			oi.variant_id,
			v.name as variant_name,
			COALESCE(
				(SELECT path FROM product_image WHERE product_id = p.id AND is_primary = true LIMIT 1),
				(SELECT path FROM product_image WHERE product_id = p.id ORDER BY id LIMIT 1)
			) as image
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN sizes s ON oi.size_id = s.id
		LEFT JOIN variants v ON oi.variant_id = v.id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		log.Printf("[GetOrderDetails] Query failed for orderID %d: %v", orderID, err)
		return nil, fmt.Errorf("failed to get order details: %w", err)
	}
	defer rows.Close()

	log.Printf("[GetOrderDetails] Query executed successfully for orderID: %d", orderID)

	var details []models.OrderDetailResponse
	rowCount := 0
	for rows.Next() {
		rowCount++
		var detail models.OrderDetailResponse
		err := rows.Scan(
			&detail.ID,
			&detail.ProductID,
			&detail.ProductName,
			&detail.Quantity,
			&detail.Price,
			&detail.SizeID,
			&detail.SizeName,
			&detail.VariantID,
			&detail.VariantName,
			&detail.Image,
		)
		if err != nil {
			log.Printf("[GetOrderDetails] Scan error on row %d for orderID %d: %v", rowCount, orderID, err)
			return nil, fmt.Errorf("failed to scan order detail: %w", err)
		}
		log.Printf("[GetOrderDetails] Row %d scanned - ProductID: %d, ProductName: %s, Quantity: %d, Image: %v",
			rowCount, detail.ProductID, detail.ProductName, detail.Quantity, detail.Image)
		details = append(details, detail)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[GetOrderDetails] Iterator error for orderID %d: %v", orderID, err)
		return nil, fmt.Errorf("error iterating order details: %w", err)
	}

	log.Printf("[GetOrderDetails] Completed for orderID: %d - Total rows scanned: %d", orderID, rowCount)

	return details, nil
}


func (p *ProductRepository) GetProductsWithSalesMetrics(ctx context.Context) ([]models.ProductSalesMetrics, error) {
	rows, err := p.db.Query(ctx, `
		SELECT 
		    pr.id,
			pr.product_name,
			COALESCE(SUM(pr.base_price * oi.quantity), 0) as revenue,
			COALESCE(SUM(oi.quantity), 0) as total_quantity
		FROM products pr
		LEFT JOIN order_items oi ON pr.id = oi.product_id
		GROUP BY pr.id, pr.product_name
		ORDER BY revenue DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query products with sales metrics: %w", err)
	}
	defer rows.Close()

	var result []models.ProductSalesMetrics

	for rows.Next() {
		var (
			productId   int64
			productName string
			revenue     int64
			quantity    int64
		)

		err := rows.Scan(
			&productId,
			&productName,
			&revenue,
			&quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, models.ProductSalesMetrics{
			ProductId  : productId,
			ProductName: productName,
			Revenue:     revenue,
			Quantity:    quantity,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}


func (p *ProductRepository) GetTopProductsByRevenue(ctx context.Context, limit int) ([]models.ProductSalesMetrics, error) {
	rows, err := p.db.Query(ctx, `
		SELECT 
			pr.product_name,
			COALESCE(SUM(pr.base_price * oi.quantity), 0) as revenue,
			COALESCE(SUM(oi.quantity), 0) as total_quantity
		FROM products pr
		LEFT JOIN order_items oi ON pr.id = oi.product_id
		GROUP BY pr.id, pr.product_name
		ORDER BY revenue DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top products by revenue: %w", err)
	}
	defer rows.Close()

	var result []models.ProductSalesMetrics

	for rows.Next() {
		var (
			productName string
			revenue     int64
			quantity    int64
		)

		err := rows.Scan(
			&productName,
			&revenue,
			&quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, models.ProductSalesMetrics{
			ProductName: productName,
			Revenue:     revenue,
			Quantity:    quantity,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}


func (p *ProductRepository) GetTopProductsByQuantity(ctx context.Context, limit int) ([]models.ProductSalesMetrics, error) {
	rows, err := p.db.Query(ctx, `
		SELECT 
			pr.product_name,
			COALESCE(SUM(pr.base_price * oi.quantity), 0) as revenue,
			COALESCE(SUM(oi.quantity), 0) as total_quantity
		FROM products pr
		LEFT JOIN order_items oi ON pr.id = oi.product_id
		GROUP BY pr.id, pr.product_name
		ORDER BY total_quantity DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top products by quantity: %w", err)
	}
	defer rows.Close()

	var result []models.ProductSalesMetrics

	for rows.Next() {
		var (
			productName string
			revenue     int64
			quantity    int64
		)

		err := rows.Scan(
			&productName,
			&revenue,
			&quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, models.ProductSalesMetrics{
			ProductName: productName,
			Revenue:     revenue,
			Quantity:    quantity,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

func (p *ProductRepository) InvalidateCache(ctx context.Context, productID int) error {
	// ✅ Correct: Delete langsung return error
	if err := p.cache.Delete(ctx, fmt.Sprintf("product:%d", productID)); err != nil {
		return fmt.Errorf("failed to invalidate product cache: %w", err)
	}

	if err := p.cache.Delete(ctx, "products:all"); err != nil {
		return fmt.Errorf("failed to invalidate products:all cache: %w", err)
	}

	return nil
}

func (p *ProductRepository) InvalidateAllProducts(ctx context.Context) error {
	return p.cache.Delete(ctx, "products:all")
}