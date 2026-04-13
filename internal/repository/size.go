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


func (r *SizeRepository) CreateSize(ctx context.Context, req *models.Size) error {
    err := r.db.QueryRow(ctx, `insert into sizes (name, additional_price) values ($1, $2) returning id`, req.Name, req.AdditionalPrice).Scan(&req.ID)
	if err != nil {
     return fmt.Errorf("failed to create sizes: %w", err)
	}
	return nil
}

func (r *SizeRepository) GetAll (ctx context.Context) ([]models.Size, error){
	 
     sizesRow, err := r.db.Query(ctx, `select id, name, additional_price from sizes`)
	 if err != nil {
		return nil, fmt.Errorf("cannot get sizes")
	 }
	 defer sizesRow.Close()

	 sizes := []models.Size{}
	 sizeIDs := []int{}

	 for sizesRow.Next() {
		var size models.Size
		err := sizesRow.Scan(
			&size.ID,
			&size.Name,
			&size.AdditionalPrice,
		
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

        sizes = append(sizes, size)
		sizeIDs = append(sizeIDs, size.ID)
	}
	 return sizes, nil
}