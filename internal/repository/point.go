package repository

import (
	"context"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PointRepository interface {
	AddPoints(ctx context.Context, tx pgx.Tx, userID int, orderID *int, points int, description string) error
	GetPointHistory(ctx context.Context, userID int, limit, offset int) ([]models.PointLedger, error)
	GetUserBalance(ctx context.Context, userID int) (int, error)
}

type pointRepository struct {
	db *pgxpool.Pool
}

func NewPointRepository(db *pgxpool.Pool) PointRepository {
	return &pointRepository{db: db}
}

func (r *pointRepository) AddPoints(ctx context.Context, tx pgx.Tx, userID int, orderID *int, points int, description string) error {
	// 1. Insert ledger
	ledgerQuery := `
		INSERT INTO point_ledgers (user_id, order_id, points, description, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`
	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, ledgerQuery, userID, orderID, points, description)
	} else {
		_, err = r.db.Exec(ctx, ledgerQuery, userID, orderID, points, description)
	}
	if err != nil {
		return err
	}

	// 2. Update users points_balance
	updateQuery := `
		UPDATE users 
		SET points_balance = points_balance + $1 
		WHERE id = $2
	`
	if tx != nil {
		_, err = tx.Exec(ctx, updateQuery, points, userID)
	} else {
		_, err = r.db.Exec(ctx, updateQuery, points, userID)
	}
	return err
}

func (r *pointRepository) GetPointHistory(ctx context.Context, userID int, limit, offset int) ([]models.PointLedger, error) {
	query := `
		SELECT id, user_id, order_id, points, description, created_at
		FROM point_ledgers
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ledgers []models.PointLedger
	for rows.Next() {
		var l models.PointLedger
		err := rows.Scan(&l.ID, &l.UserID, &l.OrderID, &l.Points, &l.Description, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		ledgers = append(ledgers, l)
	}
	return ledgers, nil
}

func (r *pointRepository) GetUserBalance(ctx context.Context, userID int) (int, error) {
	query := `SELECT points_balance FROM users WHERE id = $1`
	var balance int
	err := r.db.QueryRow(ctx, query, userID).Scan(&balance)
	return balance, err
}
