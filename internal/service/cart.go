package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) AddToCart(ctx context.Context, customerID int, req models.AddToCartRequest) (*models.CartResponse, error) {
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if err := s.cartRepo.ValidateProductSizeVariant(ctx, req.ProductID, req.SizeID, req.VariantID); err != nil {
		return nil, err
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock. available: %d, requested: %d", product.Stock, req.Quantity)
	}

	cart := &models.Cart{
		CustomerID: customerID,
		ProductID:  req.ProductID,
		SizeID:     req.SizeID,
		VariantID:  req.VariantID,
		Quantity:   req.Quantity,
	}

	if err := s.cartRepo.AddToCart(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to add to cart: %w", err)
	}

	return s.GetCart(ctx, customerID)
}

func (s *CartService) GetCart(ctx context.Context, customerID int) (*models.CartResponse, error) {
	items, err := s.cartRepo.GetCartItems(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	subtotal := 0.0
	for _, item := range items {
		subtotal += item.Subtotal
	}

	return &models.CartResponse{
		Items:      items,
		TotalItems: len(items),
		Subtotal:   subtotal,
	}, nil
}

func (s *CartService) UpdateCartItemQuantity(ctx context.Context, customerID int, cartItemID int, req models.UpdateCartItemRequest) (*models.CartResponse, error) {
	cartItem, err := s.cartRepo.GetCartItemByID(ctx, cartItemID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found: %w", err)
	}

	if cartItem.CustomerID != customerID {
		return nil, fmt.Errorf("unauthorized access to cart item")
	}

	product, err := s.productRepo.GetByID(ctx, cartItem.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock. available: %d, requested: %d", product.Stock, req.Quantity)
	}

	if err := s.cartRepo.UpdateCartQuantity(ctx, cartItemID, req.Quantity, false); err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}

	return s.GetCart(ctx, customerID)
}

func (s *CartService) RemoveFromCart(ctx context.Context, customerID int, cartItemID int) (*models.CartResponse, error) {
	cartItem, err := s.cartRepo.GetCartItemByID(ctx, cartItemID)
	if err != nil {
		return nil, fmt.Errorf("cart item not found: %w", err)
	}

	if cartItem.CustomerID != customerID {
		return nil, fmt.Errorf("unauthorized access to cart item")
	}

	if err := s.cartRepo.RemoveFromCart(ctx, cartItemID); err != nil {
		return nil, fmt.Errorf("failed to remove from cart: %w", err)
	}

	return s.GetCart(ctx, customerID)
}

func (s *CartService) ClearCart(ctx context.Context, customerID int) error {
	return s.cartRepo.ClearCart(ctx, customerID)
}

//func (s *CartService) GetCartSummary(ctx context.Context, customerID int) (map[string]interface{}, error) {
//	items, err := s.cartRepo.GetCartItems(ctx, customerID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get cart: %w", err)
//	}
//
//	subtotal := 0.0
//	totalQuantity := 0
//	for _, item := range items {
//		subtotal += item.Subtotal
//		totalQuantity += item.Quantity
//	}
//
//	return map[string]interface{}{
//		"total_items":    len(items),
//		"total_quantity": totalQuantity,
//		"subtotal":       subtotal,
//	}, nil
//}

// validateAddToCartRequest validates the add to cart request parameters
func validateAddToCartRequest(req models.AddToCartRequest) error {
	if req.ProductID <= 0 {
		return fmt.Errorf("invalid product ID")
	}

	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	if req.Quantity > 100 {
		return fmt.Errorf("quantity exceeds maximum allowed: %d", req.Quantity)
	}

	return nil
}

// executeAddToCartTransaction handles the add to cart logic with transaction support
// This method should either:
// 1. Merge with existing cart item if same product/size/temperature combo exists
func (s *CartService) executeAddToCartTransaction(ctx context.Context, customerID int, req models.AddToCartRequest) error {
	// Check if item already exists in cart
	existingItems, err := s.cartRepo.GetCartItems(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to check existing cart items: %w", err)
	}

	// Find matching cart item (same product, size, temperature)
	var existingItem *models.CartItem
	for _, item := range existingItems {
		if item.ProductID == req.ProductID &&
			item.SizeID == req.SizeID &&
			item.VariantID == req.VariantID {
			existingItem = &item
			break
		}
	}

	// If item exists, merge quantities
	if existingItem != nil {
		newQuantity := existingItem.Quantity + req.Quantity

		// Verify new total doesn't exceed stock
		product, err := s.productRepo.GetByID(ctx, req.ProductID)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		if product.Stock < newQuantity {
			return fmt.Errorf("insufficient stock for total quantity: available=%d, requested=%d", product.Stock, newQuantity)
		}

		// Update existing item - pass the cart item ID and new quantity
		// isIncrement=false means we're setting the quantity directly, not incrementing
		return s.cartRepo.UpdateCartQuantity(ctx, existingItem.ID, newQuantity, false)
	}

	// Create new cart item
	cart := &models.Cart{
		CustomerID: customerID,
		ProductID:  req.ProductID,
		SizeID:     req.SizeID,
		VariantID:  req.VariantID,
		Quantity:   req.Quantity,
	}

	return s.cartRepo.AddToCart(ctx, cart)
}

// AddToCartOptions extends AddToCart with more control over behavior
type AddToCartOptions struct {
	MergeWithExisting bool // If true, merge with existing cart item; if false, replace
	ReserveStock      bool // If true, attempt to reserve stock immediately
	ValidateOnly      bool // If true, validate but don't add to cart
}

// AddToCartWithOptions adds to cart with configurable behavior
func (s *CartService) AddToCartWithOptions(ctx context.Context, customerID int, req models.AddToCartRequest, opts *AddToCartOptions) (*models.CartResponse, error) {
	if opts == nil {
		opts = &AddToCartOptions{MergeWithExisting: true}
	}

	// Validate request
	if err := validateAddToCartRequest(req); err != nil {
		return nil, err
	}

	// Fetch and validate product
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: available=%d, requested=%d", product.Stock, req.Quantity)
	}

	if opts.ValidateOnly {
		return nil, nil
	}

	if err := s.executeAddToCartTransaction(ctx, customerID, req); err != nil {
		return nil, fmt.Errorf("failed to add to cart: %w", err)
	}

	return s.GetCart(ctx, customerID)
}
