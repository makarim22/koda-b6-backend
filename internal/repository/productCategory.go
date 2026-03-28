package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type ProductCategoryRepository struct {
	db *pgx.Conn
}

func NewProductCategoryRepository(db *pgx.Conn) *ProductCategoryRepository {
	return &ProductCategoryRepository{
		db: db,
	}
}

func (r *ProductCategoryRepository) Create(ctx context.Context, category *models.ProductCategory) error {
	query := `INSERT INTO product_category (name, description) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *ProductCategoryRepository) GetByID(ctx context.Context, id int) (*models.ProductCategory, error) {
	query := `SELECT id, name, description FROM product_category WHERE id = $1`
	category := &models.ProductCategory{}
	err := r.db.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}

func (r *ProductCategoryRepository) GetAll(ctx context.Context) ([]models.ProductCategory, error) {
	query := `SELECT id, name, description FROM product_category`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	defer rows.Close()

	var categories []models.ProductCategory
	for rows.Next() {
		var category models.ProductCategory
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

func (r *ProductCategoryRepository) Update(ctx context.Context, category *models.ProductCategory) error {
	query := `UPDATE product_category SET name=$1, description=$2 WHERE id=$3`
	commandTag, err := r.db.Exec(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("category with ID %d not found", category.ID)
	}

	return nil
}

func (r *ProductCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM product_category WHERE id=$1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("category with ID %d not found", id)
	}

	return nil
}
