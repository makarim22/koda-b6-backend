package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db *pgx.Conn
}

func NewOrderRepository(db *pgx.Conn) *OrderRepository {
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
			INSERT INTO order_detail (order_id, product_id, size_id, temperature_id, unit_price, quantity)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err := r.db.Exec(ctx, query,
			orderID,
			detail.ProductID,
			detail.SizeID,
			detail.TemperatureID,
			detail.Price,
			detail.Quantity,
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

func (r *OrderRepository) GetOrderDetails(ctx context.Context, orderID int) ([]models.OrderDetail, error) {
	query := `
		SELECT 
			od.id,
			od.order_id,
			od.product_id,
			od.size_id,
			od.temperature_id,
			od.quantity,
			p.base_price,
			p.product_name,
			s.name as size_name,
			t.label as temperature_label
		FROM order_detail od
		JOIN products p ON od.product_id = p.id
		LEFT JOIN sizes s ON od.size_id = s.id
		LEFT JOIN temperature t ON od.temperature_id = t.id
		WHERE od.order_id = $1
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order details: %w", err)
	}
	defer rows.Close()

	var details []models.OrderDetail
	for rows.Next() {
		var detail models.OrderDetail
		err := rows.Scan(
			&detail.ID,
			&detail.OrderID,
			&detail.ProductID,
			&detail.SizeID,
			&detail.TemperatureID,
			&detail.Quantity,
			&detail.Price,
			&detail.ProductName,
			&detail.SizeName,
			&detail.TemperatureLabel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order detail: %w", err)
		}
		details = append(details, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order details: %w", err)
	}

	return details, nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, customerID int, limit, offset int) ([]models.Order, error) {
	query := `
		SELECT id, customer_id, order_date, subtotal, tax, delivery_fee, status, created_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
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
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
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
			c.temperature_id,
			c.quantity,
			p.base_price,
			p.product_name,
			s.name as size_name,
			t.label as temperature_label
		FROM cart c
		JOIN products p ON c.product_id = p.id
		LEFT JOIN sizes s ON c.size_id = s.id
		LEFT JOIN temperature t ON c.temperature_id = t.id
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
			&item.TemperatureID,
			&item.Quantity,
			&item.Price,
			&item.ProductName,
			&item.SizeName,
			&item.TemperatureLabel,
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
