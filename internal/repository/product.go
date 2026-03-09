
package repository

import (
	"context"
	"errors"
	"fmt"

	"koda-b6-backend/internal/models"
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	rows, err := p.db.Query(ctx,
		`SELECT id, name, description, stock, variant_id, size_id FROM products ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all products: %w", err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Product])
	if err != nil {
		return nil, fmt.Errorf("failed to collect products: %w", err)
	}

	return products, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product

	err := p.db.QueryRow(ctx,
		`SELECT id, name, description, stock, variant_id, size_id FROM products WHERE id = $1`,
		id).Scan(&product.ID, &product.Name, &product.Description, &product.Stock, &product.VariantId, &product.SizeId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return &product, nil
}

func (p *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	err := p.db.QueryRow(ctx,
		`INSERT INTO products (name, description, stock, variant_id, size_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		product.Name, product.Description, product.Stock, product.VariantId, product.SizeId).
		Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (p *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	err := p.db.QueryRow(ctx,
		`UPDATE products 
		 SET name = $1, description = $2, stock = $3, variant_id = $4, size_id = $5
		 WHERE id = $6
		 RETURNING id, name, description, stock, variant_id, size_id`,
		product.Name, product.Description, product.Stock, product.VariantId, product.SizeId, product.ID).
		Scan(&product.ID, &product.Name, &product.Description, &product.Stock, &product.VariantId, &product.SizeId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("product with ID %d not found", product.ID)
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (p *ProductRepository) Delete(ctx context.Context, id int) error {
	commandTag, err := p.db.Exec(ctx,
		`DELETE FROM products WHERE id = $1`,
		id)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}
