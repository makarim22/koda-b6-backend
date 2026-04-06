package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) AddToCart(ctx context.Context, cart *models.Cart) error {
	var existingID int
	checkQuery := `
		SELECT id FROM cart
		WHERE customer_id = $1 AND product_id = $2 AND size_id IS NOT DISTINCT FROM $3 AND variant_id IS NOT DISTINCT FROM $4
	`

	err := r.db.QueryRow(ctx, checkQuery, cart.CustomerID, cart.ProductID, cart.SizeID, cart.VariantID).Scan(&existingID)

	if err == nil {
		return r.UpdateCartQuantity(ctx, existingID, cart.Quantity, true)
	} else if err == pgx.ErrNoRows {
		query := `
			INSERT INTO cart (customer_id, product_id, size_id, variant_id, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`

		_, err := r.db.Exec(ctx, query,
			cart.CustomerID,
			cart.ProductID,
			cart.SizeID,
			cart.VariantID,
			cart.Quantity,
		)

		if err != nil {
			return fmt.Errorf("failed to add to cart: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("failed to check cart item: %w", err)
	}
}

func (r *CartRepository) GetCartItems(ctx context.Context, customerID int) ([]models.CartItem, error) {
	query := `
		SELECT 
			c.id,
			c.product_id,
			p.product_name,
			p.base_price,
			c.quantity,
			c.size_id,
			s.name as size_name,
			c.variant_id,
			v.name as variant_name,
			c.created_at
		FROM cart c
		JOIN products p ON c.product_id = p.id
		LEFT JOIN sizes s ON c.size_id = s.id
		LEFT JOIN variants v ON c.variant_id = v.id
		WHERE c.customer_id = $1
		ORDER BY c.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.Price,
			&item.Quantity,
			&item.SizeID,
			&item.SizeName,
			&item.VariantID,
			&item.VariantName,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		item.Subtotal = item.Price * float64(item.Quantity)
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	return items, nil
}

func (r *CartRepository) UpdateCartQuantity(ctx context.Context, cartItemID int, quantity int, isIncrement bool) error {
	var query string
	if isIncrement {
		query = `
			UPDATE cart
			SET quantity = quantity + $1
			WHERE id = $2
		`
	} else {
		query = `
			UPDATE cart
			SET quantity = $1
			WHERE id = $2
		`
	}

	result, err := r.db.Exec(ctx, query, quantity, cartItemID)
	if err != nil {
		return fmt.Errorf("failed to update cart quantity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

func (r *CartRepository) RemoveFromCart(ctx context.Context, cartItemID int) error {
	query := `DELETE FROM cart WHERE id = $1`

	result, err := r.db.Exec(ctx, query, cartItemID)
	if err != nil {
		return fmt.Errorf("failed to remove from cart: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

func (r *CartRepository) ClearCart(ctx context.Context, customerID int) error {
	query := `DELETE FROM cart WHERE customer_id = $1`

	_, err := r.db.Exec(ctx, query, customerID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

func (r *CartRepository) GetCartItemByID(ctx context.Context, cartItemID int) (*models.Cart, error) {
	var cart models.Cart

	query := `
		SELECT id, customer_id, product_id, size_id, variant_id, quantity, created_at
		FROM cart
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, cartItemID).Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.ProductID,
		&cart.SizeID,
		&cart.VariantID,
		&cart.Quantity,
		&cart.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	return &cart, nil
}

func (r *CartRepository) ValidateProductSizeVariant(ctx context.Context, productID int, sizeID *int, variantID *int) error {
	// Check if product exists
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check product: %w", err)
	}
	if !exists {
		return fmt.Errorf("product %d does not exist", productID)
	}

	// Validate size if provided
	if sizeID != nil {
		err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM product_sizes WHERE size_id = $1 AND product_id = $2)", sizeID, productID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check size: %w", err)
		}
		if !exists {
			return fmt.Errorf("size %d does not exist for product %d", *sizeID, productID)
		}
	}

	// Validate variant if provided
	if variantID != nil {
		err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM product_variant WHERE variant_id = $1 AND product_id = $2)", variantID, productID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check variant: %w", err)
		}
		if !exists {
			return fmt.Errorf("variant %d does not exist for product %d", *variantID, productID)
		}
	}

	return nil
}
