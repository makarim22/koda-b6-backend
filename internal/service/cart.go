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

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock. available: %d, requested: %d", product.Stock, req.Quantity)
	}

	cart := &models.Cart{
		CustomerID:    customerID,
		ProductID:     req.ProductID,
		SizeID:        req.SizeID,
		TemperatureID: req.TemperatureID,
		Quantity:      req.Quantity,
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
