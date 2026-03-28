package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type OrderDetailRepository struct {
	db *pgx.Conn
}

func NewOrderDetailRepository(db *pgx.Conn) *OrderDetailRepository {
	return &OrderDetailRepository{
		db: db,
	}
}

func (r *OrderDetailRepository) Create(ctx context.Context, detail *models.OrderDetail) error {
	query := `INSERT INTO order_detail (order_id, product_id, size_id, temperature_id, quantity, unit_price) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(ctx, query, detail.OrderID, detail.ProductID, detail.SizeID, detail.TemperatureID, detail.Quantity, detail.Price).
		Scan(&detail.ID)
	if err != nil {
		return fmt.Errorf("failed to create order detail: %w", err)
	}
	return nil
}

func (r *OrderDetailRepository) GetByOrderID(ctx context.Context, orderID int) ([]models.OrderDetail, error) {
	query := `SELECT id, order_id, product_id, size_id, temperature_id, quantity FROM order_detail WHERE order_id = $1`
	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}
	defer rows.Close()

	var details []models.OrderDetail
	for rows.Next() {
		var detail models.OrderDetail
		err := rows.Scan(&detail.ID, &detail.OrderID, &detail.ProductID, &detail.SizeID, &detail.TemperatureID, &detail.Quantity)
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
