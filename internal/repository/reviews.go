package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (r *ReviewsRepository) GetAll(ctx context.Context) ([]models.ReviewsResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT ur.id, ur.user_id, u.full_name, u.email, ur.product_id, p.product_name, ur.order_id, ur.message, ur.rating FROM user_review ur join products p on p.id = ur.product_id join users u on u.id = ur.user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ReviewsResponse])
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewsRepository) GetById(ctx context.Context, id int) (*models.ReviewsResponse, error) {
	row, err := r.db.Query(ctx,
		`SELECT ur.id, ur.user_id, u.full_name, u.email, ur.product_id, p.product_name, ur.order_id, ur.message, ur.rating 
		FROM user_review ur 
		JOIN products p ON p.id = ur.product_id 
		JOIN users u ON u.id = ur.user_id 
		WHERE ur.id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	review, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReviewsResponse])
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewsRepository) CreateReview(ctx context.Context, review *models.ReviewsRequest) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO user_review (user_id, product_id, order_id, message, rating) 
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id`,
		review.UserId, review.ProductId, review.OrderId, review.Message, review.Rating).Scan(&review.Id)
	return err
}

func (r *ReviewsRepository) GetByUserProductOrder(ctx context.Context, userID, productID, orderID int) (*models.Reviews, error) {
	var review models.Reviews

	err := r.db.QueryRow(ctx,
		`SELECT id FROM user_review WHERE user_id = $1 AND product_id = $2 AND order_id = $3`,
		userID, productID, orderID).Scan(&review.Id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check review: %w", err)
	}

	return &review, nil
}

func (r *ReviewsRepository) UpdateReview(ctx context.Context, review *models.ReviewsRequest) error {
	_, err := r.db.Exec(ctx,
		`UPDATE user_review SET message=$1, rating=$2 where id=$3`,
		review.Message, review.Rating, review.Id)
	return err
}

func (r *ReviewsRepository) GetAverageRating(ctx context.Context, productID int) (float64, error) {
	query := `SELECT AVG(rating) FROM user_review WHERE product_id = $1`
	var rating sql.NullFloat64
	err := r.db.QueryRow(ctx, query, productID).Scan(&rating)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}
	if !rating.Valid {
		return 0, nil
	}
	return rating.Float64, nil
}
