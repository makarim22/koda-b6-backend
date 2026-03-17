package repository

import (
	"context"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type ReviewsRepository struct {
	db *pgx.Conn
}

func NewReviewsRepository(db *pgx.Conn) *ReviewsRepository {
	return &ReviewsRepository{
		db: db,
	}
}

func (r *ReviewsRepository) GetAll(ctx context.Context) ([]models.Reviews, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, product_id, order_id, message, rating FROM user_review`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Reviews])
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
