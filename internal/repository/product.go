package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	rows, err := p.db.Query(ctx, `
		SELECT 
			pr.id,
			pr.product_name,
			pr.description,
			pr.stock,
			pr.base_price,
			pi.id as image_id,
			pi.product_id,
			pi.path,
			pi.is_primary
		FROM products pr
		LEFT JOIN product_image pi ON pr.id = pi.product_id
		ORDER BY pr.id DESC, pi.is_primary DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all products: %w", err)
	}
	defer rows.Close()

	// Map to collect products with their images
	productMap := make(map[int]*models.Product)
	productOrder := []int{}

	for rows.Next() {
		var (
			id          int
			productName string
			description string
			stock       int
			basePrice   int
			imageID     *int
			productID   *int
			imagePath   *string
			isPrimary   *bool
		)

		err := rows.Scan(
			&id,
			&productName,
			&description,
			&stock,
			&basePrice,
			&imageID,
			&productID,
			&imagePath,
			&isPrimary,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}

		// If product not yet in map, create it
		if _, exists := productMap[id]; !exists {
			productMap[id] = &models.Product{
				ID:          id,
				ProductName: productName,
				Description: description,
				Stock:       stock,
				BasePrice:   basePrice,
				Images:      []models.ProductImage{},
			}
			productOrder = append(productOrder, id)
		}

		// Add image if it exists
		if imageID != nil && imagePath != nil {
			image := models.ProductImage{
				ID:        *imageID,
				ProductID: *productID,
				Path:      *imagePath,
				IsPrimary: *isPrimary,
			}
			productMap[id].Images = append(productMap[id].Images, image)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	// Convert map back to slice in original order
	products := make([]models.Product, len(productOrder))
	for i, id := range productOrder {
		products[i] = *productMap[id]
	}

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
	err := p.db.QueryRow(ctx,
		`INSERT INTO products (product_name, description, stock, base_price)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		product.ProductName, product.Description, product.Stock, product.BasePrice).
		Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
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
			v.name as variant_name
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
		)
		if err != nil {
			log.Printf("[GetOrderDetails] Scan error on row %d for orderID %d: %v", rowCount, orderID, err)
			return nil, fmt.Errorf("failed to scan order detail: %w", err)
		}
		log.Printf("[GetOrderDetails] Row %d scanned - ProductID: %d, ProductName: %s, Quantity: %d",
			rowCount, detail.ProductID, detail.ProductName, detail.Quantity)
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
