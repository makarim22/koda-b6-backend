package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductDiscountRepository struct {
	db *pgxpool.Pool
}

func NewProductDiscountRepository(db *pgxpool.Pool) *ProductDiscountRepository {
	return &ProductDiscountRepository{
		db: db,
	}
}

func (r *ProductDiscountRepository) GetByProductID(ctx context.Context, productID int) ([]models.ProductDiscount, error) {
	query := `SELECT id, product_id, discount_rate, description, is_flash_sale FROM product_discount WHERE product_id = $1`
	rows, err := r.db.Query(ctx, query, productID)

	if err != nil {
		return nil, fmt.Errorf("failed to query discounts for product %d: %w", productID, err)
	}
	defer rows.Close()

	var discounts []models.ProductDiscount
	for rows.Next() {
		var discount models.ProductDiscount

		if err := rows.Scan(&discount.ID, &discount.ProductID, &discount.DiscountRate, &discount.Description, &discount.IsFlashSale); err != nil {
			return nil, fmt.Errorf("failed to scan discount row: %w", err)
		}
		discounts = append(discounts, discount)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate discounts: %w", err)
	}
	return discounts, nil
}
