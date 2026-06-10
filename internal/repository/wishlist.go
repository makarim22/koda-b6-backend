package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WishlistRepository struct {
	pool *pgxpool.Pool
}

func NewWishlistRepository(pool *pgxpool.Pool) *WishlistRepository {
	return &WishlistRepository{pool: pool}
}

func (r *WishlistRepository) Add(ctx context.Context, customerID int, productID int) error {
	query := `
		INSERT INTO wishlists (customer_id, product_id)
		VALUES ($1, $2)
		ON CONFLICT (customer_id, product_id) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query, customerID, productID)
	if err != nil {
		return fmt.Errorf("failed to add to wishlist: %w", err)
	}
	return nil
}

func (r *WishlistRepository) Remove(ctx context.Context, customerID int, productID int) error {
	query := `DELETE FROM wishlists WHERE customer_id = $1 AND product_id = $2`
	_, err := r.pool.Exec(ctx, query, customerID, productID)
	if err != nil {
		return fmt.Errorf("failed to remove from wishlist: %w", err)
	}
	return nil
}

func (r *WishlistRepository) GetUserWishlist(ctx context.Context, customerID int) ([]models.WishlistItemResponse, error) {
	query := `
		SELECT 
			w.product_id, 
			p.product_name, 
			p.base_price, 
			p.description,
			w.created_at as added_at,
			(SELECT path FROM product_image pi WHERE pi.product_id = p.id AND pi.is_primary = true LIMIT 1) as image
		FROM wishlists w
		JOIN products p ON w.product_id = p.id
		WHERE w.customer_id = $1
		ORDER BY w.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}
	defer rows.Close()

	var items []models.WishlistItemResponse
	for rows.Next() {
		var item models.WishlistItemResponse
		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.BasePrice,
			&item.Description,
			&item.AddedAt,
			&item.Image,
		); err != nil {
			return nil, fmt.Errorf("failed to scan wishlist item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *WishlistRepository) CheckStatus(ctx context.Context, customerID int, productID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM wishlists WHERE customer_id = $1 AND product_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, customerID, productID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check wishlist status: %w", err)
	}
	return exists, nil
}
