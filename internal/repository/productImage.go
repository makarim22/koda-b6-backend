package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *ProductImageRepository) Save(ctx context.Context, image *models.ProductImage) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO product_image (product_id, path, is_primary)
         VALUES ($1, $2, $3)
         RETURNING id`,
		image.ProductID, image.Path, image.IsPrimary).
		Scan(&image.ID)

	if err != nil {
		return fmt.Errorf("failed to save product image: %w", err)
	}
	return nil
}

func (r *ProductImageRepository) CountByProductID(ctx context.Context, productID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM product_image WHERE product_id = $1`,
		productID).
		Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count images: %w", err)
	}
	return count, nil
}

func (r *ProductImageRepository) UnsetPrimary(ctx context.Context, productID int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE product_image SET is_primary = false
         WHERE product_id = $1 AND is_primary = true`,
		productID)

	if err != nil {
		return fmt.Errorf("failed to unset primary image: %w", err)
	}
	return nil
}

func (r *ProductImageRepository) SetPrimary(ctx context.Context, imageID, productID int) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE product_image SET is_primary = true
         WHERE id = $1 AND product_id = $2`,
		imageID, productID)

	if err != nil {
		return fmt.Errorf("failed to set primary image: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("image with ID %d not found for product %d", imageID, productID)
	}
	return nil
}

func (r *ProductImageRepository) Delete(ctx context.Context, imageID, productID int) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM product_image WHERE id = $1 AND product_id = $2`,
		imageID, productID)

	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("image with ID %d not found for product %d", imageID, productID)
	}
	return nil
}

func (r *ProductImageRepository) FindByID(ctx context.Context, imageID, productID int) (*models.ProductImage, error) {
	var image models.ProductImage

	err := r.db.QueryRow(ctx,
		`SELECT id, product_id, path, is_primary
         FROM product_image
         WHERE id = $1 AND product_id = $2`,
		imageID, productID).
		Scan(&image.ID, &image.ProductID, &image.Path, &image.IsPrimary)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find image by id: %w", err)
	}
	return &image, nil
}

func (r *ProductImageRepository) PromoteNextPrimary(ctx context.Context, productID int) error {
	// Pick the lowest id image and promote it
	tag, err := r.db.Exec(ctx,
		`UPDATE product_image SET is_primary = true
				 WHERE id = (
					 SELECT id FROM product_image
					 WHERE product_id = $1
					 ORDER BY id
					 LIMIT 1
				 )`,
		productID)

	if err != nil {
		return fmt.Errorf("failed to promote next primary image: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil // no images left, nothing to promote
	}
	return nil
}
