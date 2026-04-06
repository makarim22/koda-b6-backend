package repository

import (
	"context"
	"fmt"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"koda-b6-backend/internal/models"
)

type ProductImageRepository struct {
	db *pgxpool.Pool
}

func NewProductImageRepository(db *pgxpool.Pool) *ProductImageRepository {
	return &ProductImageRepository{
		db: db,
	}
}

func (r *ProductImageRepository) GetByProductImageID(ctx context.Context, productID int) ([]models.ProductImage, error) {
	query := `SELECT id, product_id, path, is_primary FROM product_image WHERE product_id = $1`
	rows, err := r.db.Query(ctx, query, productID)

	if err != nil {
		return nil, fmt.Errorf("failed to query images for product %d: %w", productID, err)
	}
	defer rows.Close()

	var images []models.ProductImage
	for rows.Next() {
		var image models.ProductImage

		if err := rows.Scan(&image.ID, &image.ProductID, &image.Path, &image.IsPrimary); err != nil {
			return nil, fmt.Errorf("failed to scan images row: %w", err)
		}
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate images: %w", err)
	}
	return images, nil
}
