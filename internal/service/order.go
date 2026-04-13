package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID int, req models.CreateOrderRequest) (*models.OrderResponse, error) {
	cartItems, err := s.orderRepo.GetCartItems(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	subtotal := 0.0
	for _, item := range cartItems {
		subtotal += item.Price * float64(item.Quantity)
	}

	order := &models.Order{
		CustomerID:  customerID,
		Subtotal:    subtotal,
		Tax:         req.Tax,
		DeliveryFee: req.DeliveryFee,
		Status:      "pending",
	}

	orderID, err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if err := s.orderRepo.CreateOrderDetails(ctx, orderID, cartItems); err != nil {
		return nil, fmt.Errorf("failed to create order details: %w", err)
	}

	if err := s.orderRepo.ClearCart(ctx, customerID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	return s.GetOrder(ctx, orderID, customerID)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID int, customerID int) (*models.OrderResponse, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.CustomerID != customerID {
		return nil, fmt.Errorf("unauthorized access to order")
	}

	details, err := s.orderRepo.GetOrderDetails(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order details: %w", err)
	}

	response := &models.OrderResponse{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		OrderDate:   order.OrderDate,
		Subtotal:    order.Subtotal,
		Tax:         order.Tax,
		DeliveryFee: order.DeliveryFee,
		Total:       order.Subtotal + order.Tax + order.DeliveryFee,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
	}

	for _, detail := range details {
		itemResponse := models.OrderDetailResponse{
			ID:          detail.ID,
			ProductID:   detail.ProductID,
			ProductName: detail.ProductName,
			Quantity:    detail.Quantity,
			Price:       detail.Price,
			SizeID:      detail.SizeID,
			SizeName:    detail.SizeName,
			VariantID:   detail.VariantID,
			VariantName: detail.VariantName,
		}
		response.Items = append(response.Items, itemResponse)
	}

	return response, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, customerID int, limit, offset int) ([]models.OrderResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	orders, err := s.orderRepo.GetUserOrders(ctx, customerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		response := models.OrderResponse{
			ID:          order.ID,
			CustomerID:  order.CustomerID,
			OrderDate:   order.OrderDate,
			Subtotal:    order.Subtotal,
			Tax:         order.Tax,
			DeliveryFee: order.DeliveryFee,
			Total:       order.Subtotal + order.Tax + order.DeliveryFee,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			Items:       order.Items,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int, customerID int, newStatus string) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.CustomerID != customerID {
		return fmt.Errorf("unauthorized access to order")
	}

	validTransitions := map[string][]string{
		"pending":    {"processing", "cancelled"},
		"processing": {"shipped", "cancelled"},
		"shipped":    {"delivered"},
		"delivered":  {},
		"cancelled":  {},
	}

	allowedStatuses, exists := validTransitions[order.Status]
	if exists {
		isAllowed := false
		for _, status := range allowedStatuses {
			if status == newStatus {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return fmt.Errorf("invalid status transition from %s to %s", order.Status, newStatus)
		}
	}

	return s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus)
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderID int, customerID int) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.CustomerID != customerID {
		return fmt.Errorf("unauthorized access to order")
	}

	if order.Status != "pending" {
		return fmt.Errorf("can only delete pending orders, current status: %s", order.Status)
	}

	return s.orderRepo.DeleteOrder(ctx, orderID)
}



func (s *OrderService) GetDailySales(ctx context.Context) ([]models.DailySalesData, error) {
    result, err := s.orderRepo.GetDailySalesData(ctx)
    if err != nil {
        return nil, fmt.Errorf("cannot retrieve daily sales data: %w", err)  // fixed typo + better message
    }
    return result, nil
}