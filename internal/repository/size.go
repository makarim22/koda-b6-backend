package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"koda-b6-backend/internal/models"
	//"strconv"
)

type SizeRepository struct {
	db *pgxpool.Pool
}

func NewSizeRepository(db *pgxpool.Pool) *SizeRepository {
	return &SizeRepository{
		db: db,
	}
}

func (r *SizeRepository) GetSizesByProductID(ctx context.Context, productID int) ([]models.Size, error) {
	query := `
		SELECT DISTINCT
			s.id,
			s.name,
			s.additional_price
		FROM sizes s
		INNER JOIN product_sizes ps ON s.id = ps.size_id
		WHERE ps.product_id = $1
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sizes for product %d: %w", productID, err)
	}
	defer rows.Close()

	var sizes []models.Size

	for rows.Next() {
		var size models.Size

		if err := rows.Scan(&size.ID, &size.Name, &size.AdditionalPrice); err != nil {
			return nil, fmt.Errorf("failed to scan size row: %w", err)
		}

		sizes = append(sizes, size)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating size rows: %w", err)
	}

	return sizes, nil
}
