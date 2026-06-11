package repository

import (
	"context"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackingRepository interface {
	InsertTracking(ctx context.Context, tx pgx.Tx, orderID int, status string, description *string) error
	GetTrackingByOrderID(ctx context.Context, orderID int) ([]models.OrderTracking, error)
}

type trackingRepository struct {
	db *pgxpool.Pool
}

func NewTrackingRepository(db *pgxpool.Pool) TrackingRepository {
	return &trackingRepository{db: db}
}

func (r *trackingRepository) InsertTracking(ctx context.Context, tx pgx.Tx, orderID int, status string, description *string) error {
	query := `
		INSERT INTO order_tracking (order_id, status, description, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
	`
	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, orderID, status, description)
	} else {
		_, err = r.db.Exec(ctx, query, orderID, status, description)
	}
	return err
}

func (r *trackingRepository) GetTrackingByOrderID(ctx context.Context, orderID int) ([]models.OrderTracking, error) {
	query := `
		SELECT id, order_id, status, description, created_at
		FROM order_tracking
		WHERE order_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trackings []models.OrderTracking
	for rows.Next() {
		var t models.OrderTracking
		err := rows.Scan(&t.ID, &t.OrderID, &t.Status, &t.Description, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		trackings = append(trackings, t)
	}
	return trackings, nil
}
