package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type VariantRepository struct {
	db *pgx.Conn
}

func NewVariantRepository(db *pgx.Conn) *VariantRepository {
	return &VariantRepository{
		db: db,
	}
}

func (r *VariantRepository) GetVariantsByProductID(ctx context.Context, productID int) ([]models.Variant, error) {
	query := `select v.id, v.name, v.additional_price from variants v join product_variant pv on v.id = pv.variant_id where pv.product_id = $1`
	rows, err := r.db.Query(ctx, query, productID)

	if err != nil {
		return nil, fmt.Errorf("failed to query variants for product %d: %w", productID, err)
	}
	defer rows.Close()

	var variants []models.Variant
	for rows.Next() {
		var variant models.Variant

		if err := rows.Scan(&variant.ID, &variant.Name, &variant.AdditionalPrice); err != nil {
			return nil, fmt.Errorf("failed to scan variant row: %w", err)

		}
		variants = append(variants, variant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate variants: %w", err)
	}
	return variants, nil
}
