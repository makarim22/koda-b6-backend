package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) (int, error) {
	var orderID int

	query := `
		INSERT INTO orders (customer_id, subtotal, tax, delivery_fee, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		order.CustomerID,
		order.Subtotal,
		order.Tax,
		order.DeliveryFee,
		order.Status,
	).Scan(&orderID)

	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	return orderID, nil
}

func (r *OrderRepository) CreateOrderDetails(ctx context.Context, orderID int, details []models.OrderDetail) error {
	for _, detail := range details {
		query := `
			INSERT INTO order_items (order_id, product_id, size_id, variant_id, unit_price, quantity, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		_, err := r.db.Exec(ctx, query,
			orderID,
			detail.ProductID,
			detail.SizeID,
			detail.VariantID,
			detail.Price,
			detail.Quantity,
			detail.Subtotal,
		)

		if err != nil {
			return fmt.Errorf("failed to create order detail: %w", err)
		}
	}

	return nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order

	query := `
		SELECT id, customer_id, order_date, subtotal, tax, delivery_fee, status, created_at
		FROM orders
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.CustomerID,
		&order.OrderDate,
		&order.Subtotal,
		&order.Tax,
		&order.DeliveryFee,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (r *OrderRepository) GetCartItems(ctx context.Context, customerID int) ([]models.OrderDetail, error) {
	query := `
		SELECT 
			c.id,
			0 as order_id,
			c.product_id,
			c.size_id,
			c.variant_id,
			c.quantity,
			p.base_price,
			p.product_name,
			s.name as size_name,
			v.name as variant_name,
		    (p.base_price * c.quantity) as subtotal
		FROM cart c
		JOIN products p ON c.product_id = p.id
		LEFT JOIN sizes s ON c.size_id = s.id
		LEFT JOIN variants v ON c.variant_id = v.id
		WHERE c.customer_id = $1
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderDetail
	for rows.Next() {
		var item models.OrderDetail
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.SizeID,
			&item.VariantID,
			&item.Quantity,
			&item.Price,
			&item.ProductName,
			&item.SizeName,
			&item.VariantName,
			&item.Subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *OrderRepository) ClearCart(ctx context.Context, customerID int) error {
	query := `DELETE FROM cart WHERE customer_id = $1`

	_, err := r.db.Exec(ctx, query, customerID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, orderID int) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, customerID int, limit, offset int) ([]models.OrderResponse, error) {
	query := `
		SELECT id, customer_id, order_date, subtotal, tax, delivery_fee, status, created_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	log.Printf("[GetUserOrders] Starting query for customerID: %d, limit: %d, offset: %d", customerID, limit, offset)

	rows, err := r.db.Query(ctx, query, customerID, limit, offset)
	if err != nil {
		log.Printf("[GetUserOrders] Query failed: %v", err)
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()

	var orders []models.OrderResponse
	orderCount := 0

	for rows.Next() {
		var order models.OrderResponse
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.OrderDate,
			&order.Subtotal,
			&order.Tax,
			&order.DeliveryFee,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			log.Printf("[GetUserOrders] Scan failed for order: %v", err)
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		log.Printf("[GetUserOrders] Scanned order %d for customerID %d", order.ID, order.CustomerID)

		// Fetch order items for this order
		items, err := r.GetOrderDetails(ctx, order.ID)
		if err != nil {
			log.Printf("[GetUserOrders] Failed to fetch items for order %d: %v", order.ID, err)
			return nil, fmt.Errorf("failed to fetch items for order %d: %w", order.ID, err)
		}

		log.Printf("[GetUserOrders] Fetched %d items for order %d, assigning to order.Items", len(items), order.ID)
		order.Items = items

		log.Printf("[GetUserOrders] Order %d - Items assigned: %d items in order.Items slice", order.ID, len(order.Items))

		orders = append(orders, order)
		orderCount++
	}

	if err = rows.Err(); err != nil {
		log.Printf("[GetUserOrders] Iterator error: %v", err)
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	log.Printf("[GetUserOrders] Completed for customerID %d - Total orders returned: %d", customerID, orderCount)
	return orders, nil
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
