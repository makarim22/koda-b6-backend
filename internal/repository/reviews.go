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
		`SELECT ur.id, ur.user_id, u.full_name, u.email, ur.product_id, p.product_name, ur.order_id, ur.message, ur.rating FROM user_review ur join products p on p.id = ur.product_id join users u on u.id = ur.user_id`)
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
