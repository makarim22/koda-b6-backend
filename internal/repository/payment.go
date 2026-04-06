package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `INSERT INTO payment (order_id, method, amount, status) 
	          VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, payment.OrderID, payment.Method, payment.Amount, "pending").
		Scan(&payment.ID)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID int) (*models.Payment, error) {
	query := `SELECT id, order_id, method, amount, status, transaction_id, payment_date FROM payment WHERE order_id = $1`
	payment := &models.Payment{}
	err := r.db.QueryRow(ctx, query, orderID).Scan(&payment.ID, &payment.OrderID, &payment.Method, &payment.Amount, &payment.Status, &payment.TransactionID, &payment.PaymentDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment for order ID %d not found", orderID)
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE payment SET status=$1 WHERE id=$2`
	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment with ID %d not found", id)
	}
	return nil
}

func (r *PaymentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM payment WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment with ID %d not found", id)
	}
	return nil
}
