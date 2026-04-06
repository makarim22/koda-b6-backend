package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderDetailRepository struct {
	db *pgxpool.Pool
}

func NewOrderDetailRepository(db *pgxpool.Pool) *OrderDetailRepository {
	return &OrderDetailRepository{
		db: db,
	}
}

func (r *OrderDetailRepository) Create(ctx context.Context, detail *models.OrderDetail) error {
	query := `INSERT INTO order_detail (order_id, product_id, size_id, temperature_id, quantity, unit_price) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(ctx, query, detail.OrderID, detail.ProductID, detail.SizeID, detail.VariantID, detail.Quantity, detail.Price).
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
		err := rows.Scan(&detail.ID, &detail.OrderID, &detail.ProductID, &detail.SizeID, &detail.VariantID, &detail.Quantity)
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

func (r *OrderDetailRepository) Update(ctx context.Context, detail *models.OrderDetail) error {
	query := `UPDATE order_detail SET quantity = $1, unit_price = $2 WHERE id = $3`
	result, err := r.db.Exec(ctx, query, detail.Quantity, detail.Price, detail.ID)
	if err != nil {
		return fmt.Errorf("failed to update order detail: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("order detail not found")
	}
	return nil
}

func (r *OrderDetailRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM order_detail WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order detail: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("order detail not found")
	}
	return nil
}
