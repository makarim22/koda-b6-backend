package repository

import (
	"context"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type ProductDiscountRepository struct {
	db *pgx.Conn
}

func NewProductDiscountRepository(db *pgx.Conn) *ProductDiscountRepository {
	return &ProductDiscountRepository{
		db: db,
	}
}

func (r *ProductDiscountRepository) GetByProductID(ctx context.Context, productID int) (*models.ProductDiscount, error) {
	var discount models.ProductDiscount

	err := r.db.QueryRow(ctx,
		`SELECT id, product_id, discount_rate, description, is_flash_sale FROM product_discount WHERE product_id = $1`,
		productID).Scan(&discount.ID, &discount.ProductID, &discount.DiscountRate, &discount.Description, &discount.IsFlashSale)

	if err != nil {
		return nil, nil
	}

	return &discount, nil
}
