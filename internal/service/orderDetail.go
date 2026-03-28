package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"

	defaultErr "errors"
	"github.com/jackc/pgx/v5"
)

type OrderDetailService struct {
	orderDetailRepo *repository.OrderDetailRepository
	orderRepo       *repository.OrderRepository
	productRepo     *repository.ProductRepository
}

func NewOrderDetailService(
	orderDetailRepo *repository.OrderDetailRepository,
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
) *OrderDetailService {
	return &OrderDetailService{
		orderDetailRepo: orderDetailRepo,
		orderRepo:       orderRepo,
		productRepo:     productRepo,
	}
}

// Create validates and creates a new order detail
func (s *OrderDetailService) Create(ctx context.Context, detail *models.OrderDetail) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	// Validate required fields
	if detail.OrderID <= 0 {
		return errors.NewValidationError("order_id", "must be greater than 0")
	}
	if detail.ProductID <= 0 {
		return errors.NewValidationError("product_id", "must be greater than 0")
	}
	if detail.Quantity <= 0 {
		return errors.NewValidationError("quantity", "must be greater than 0")
	}
	if detail.Price < 0 {
		return errors.NewValidationError("unit_price", "cannot be negative")
	}

	// Verify order exists
	_, err := s.orderRepo.GetOrderByID(ctx, detail.OrderID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return errors.NewValidationError("order_id", "does not exist")
		}
		return fmt.Errorf("failed to verify order: %w", err)
	}

	// Verify product exists
	_, err = s.productRepo.GetByID(ctx, detail.ProductID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return errors.NewValidationError("product_id", "does not exist")
		}
		return fmt.Errorf("failed to verify product: %w", err)
	}

	// Create order detail
	err = s.orderDetailRepo.Create(ctx, detail)
	if err != nil {
		return fmt.Errorf("failed to create order detail: %w", err)
	}

	return nil
}

// GetByOrderID retrieves all order details for a specific order
func (s *OrderDetailService) GetByOrderID(ctx context.Context, orderID int) ([]*models.OrderDetail, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	if orderID <= 0 {
		return nil, errors.NewValidationError("order_id", "must be greater than 0")
	}

	_, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		if defaultErr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewNotFoundError("Order", orderID)
		}
		return nil, fmt.Errorf("failed to verify order: %w", err)
	}

	// Fetch order details
	details, err := s.orderDetailRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}

	// Convert slice of models to slice of pointers
	result := make([]*models.OrderDetail, len(details))
	for i := range details {
		result[i] = &details[i]
	}

	return result, nil
}

//// GetByID retrieves a single order detail by ID
//func (s *OrderDetailService) GetByID(ctx context.Context, id int) (*models.OrderDetail, error) {
//	if err := ctx.Err(); err != nil {
//		return nil, fmt.Errorf("context error: %w", err)
//	}
//
//	if id <= 0 {
//		return nil, errors.NewValidationError("id", "must be greater than 0")
//	}
//
//	// Verify order detail exists
//	details, err := s.orderDetailRepo.GetByOrderID(ctx, id)
//	if err != nil {
//		return fmt.Errorf("failed to verify order detail: %w", err)
//	}
//	if len(details) == 0 {
//		return errors.NewNotFoundError("OrderDetail", id)
//	}
//
//	return details, nil
//}

// Update validates and updates an order detail
func (s *OrderDetailService) Update(ctx context.Context, detail *models.OrderDetail) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	if detail.ID <= 0 {
		return errors.NewValidationError("id", "must be greater than 0")
	}

	// Verify order detail exists
	details, err := s.orderDetailRepo.GetByOrderID(ctx, detail.ID)
	if err != nil {
		return fmt.Errorf("failed to verify order detail: %w", err)
	}
	if len(details) == 0 {
		return errors.NewNotFoundError("OrderDetail", detail.ID)
	}

	// Validate updateable fields
	if detail.Quantity <= 0 {
		return errors.NewValidationError("quantity", "must be greater than 0")
	}
	if detail.Price < 0 {
		return errors.NewValidationError("unit_price", "cannot be negative")
	}

	err = s.orderDetailRepo.Update(ctx, detail)
	if err != nil {
		return fmt.Errorf("failed to update order detail: %w", err)
	}

	return nil
}

// Delete removes an order detail by ID
func (s *OrderDetailService) Delete(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	if id <= 0 {
		return errors.NewValidationError("id", "must be greater than 0")
	}

	// Verify order detail exists
	details, err := s.orderDetailRepo.GetByOrderID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to verify order detail: %w", err)
	}
	if len(details) == 0 {
		return errors.NewNotFoundError("OrderDetail", id)
	}

	err = s.orderDetailRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete order detail: %w", err)
	}

	return nil
}
