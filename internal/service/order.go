package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
	"net/http"
	"os"
)

type OrderService struct {
	orderRepo      *repository.OrderRepository
	productRepo    *repository.ProductRepository
	paymentRepo    *repository.PaymentRepository
	voucherService *VoucherService
	pointService   PointService
	trackingRepo   repository.TrackingRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, paymentRepo *repository.PaymentRepository, voucherService *VoucherService, pointService PointService, trackingRepo repository.TrackingRepository) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		paymentRepo:    paymentRepo,
		voucherService: voucherService,
		pointService:   pointService,
		trackingRepo:   trackingRepo,
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

	var voucherID *int
	var discountAmount float64

	if req.VoucherCode != "" {
		voucher, discAmt, err := s.voucherService.CalculateDiscount(ctx, req.VoucherCode, subtotal)
		if err != nil {
			return nil, fmt.Errorf("voucher error: %w", err)
		}
		vid := voucher.ID
		voucherID = &vid
		discountAmount = discAmt
	}

	// Calculate Point Usage
	// 1 Point = 10 IDR
	pointDiscount := 0.0
	pointsUsed := 0
	if req.PointsToUse > 0 {
		balance, err := s.pointService.GetUserBalance(ctx, customerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user point balance: %w", err)
		}
		if balance < req.PointsToUse {
			return nil, fmt.Errorf("insufficient points balance")
		}
		
		pointsUsed = req.PointsToUse
		pointDiscount = float64(pointsUsed * 10)
		
		// Ensure discount doesn't exceed subtotal
		if pointDiscount > subtotal-discountAmount {
			pointDiscount = subtotal - discountAmount
			pointsUsed = int(pointDiscount / 10)
		}
		discountAmount += pointDiscount
	}

	order := &models.Order{
		CustomerID:     customerID,
		Subtotal:       subtotal,
		Tax:            req.Tax,
		DeliveryFee:    req.DeliveryFee,
		Status:         "pending",
		VoucherID:      voucherID,
		DiscountAmount: discountAmount,
		PointsUsed:     pointsUsed,
	}

	orderID, err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	if voucherID != nil {
		if err := s.voucherService.voucherRepo.IncrementUsage(ctx, *voucherID); err != nil {
			// Intentionally not failing the order if usage increment fails
			fmt.Printf("failed to increment voucher usage: %v\n", err)
		}
	}

	if pointsUsed > 0 {
		err := s.pointService.DeductPoints(ctx, nil, customerID, &orderID, pointsUsed, "Used points for order checkout")
		if err != nil {
			fmt.Printf("failed to deduct points: %v\n", err)
			// in a real prod system, this should be inside a TX
		}
	}

	if err := s.orderRepo.CreateOrderDetails(ctx, orderID, cartItems); err != nil {
		return nil, fmt.Errorf("failed to create order details: %w", err)
	}

	// Insert tracking log
	desc := "Order has been placed"
	_ = s.trackingRepo.InsertTracking(ctx, nil, orderID, "pending", &desc)

	if err := s.orderRepo.ClearCart(ctx, customerID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Create a pending payment record
	payment := &models.Payment{
		OrderID: orderID,
		Amount:  subtotal + req.Tax + req.DeliveryFee - discountAmount,
		Method:  "credit_card", // default or placeholder for Snap
		Status:  "pending",
	}
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		fmt.Printf("failed to create pending payment record: %v\n", err)
	}

	// Generate Midtrans Snap Token and Redirect URL
	snapToken, redirectURL, err := s.GenerateSnapToken(orderID, payment.Amount)
	if err != nil {
		fmt.Printf("failed to generate snap token: %v\n", err)
	}

	orderResp, err := s.GetOrder(ctx, orderID, customerID)
	if err != nil {
		return nil, err
	}
	orderResp.SnapToken = snapToken
	orderResp.RedirectURL = redirectURL

	return orderResp, nil
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
		ID:             order.ID,
		CustomerID:     order.CustomerID,
		OrderDate:      order.OrderDate,
		Subtotal:       order.Subtotal,
		Tax:            order.Tax,
		DeliveryFee:    order.DeliveryFee,
		DiscountAmount: order.DiscountAmount,
		Total:          order.Subtotal + order.Tax + order.DeliveryFee - order.DiscountAmount,
		Status:         order.Status,
		CreatedAt:      order.CreatedAt,
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
			ID:             order.ID,
			CustomerID:     order.CustomerID,
			OrderDate:      order.OrderDate,
			Subtotal:       order.Subtotal,
			Tax:            order.Tax,
			DeliveryFee:    order.DeliveryFee,
			DiscountAmount: order.DiscountAmount,
			Total:          order.Subtotal + order.Tax + order.DeliveryFee - order.DiscountAmount,
			Status:         order.Status,
			CreatedAt:      order.CreatedAt,
			Items:          order.Items,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, limit, offset int) ([]models.OrderResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	orders, err := s.orderRepo.GetAllOrders(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		response := models.OrderResponse{
			ID:             order.ID,
			CustomerID:     order.CustomerID,
			OrderDate:      order.OrderDate,
			Subtotal:       order.Subtotal,
			Tax:            order.Tax,
			DeliveryFee:    order.DeliveryFee,
			DiscountAmount: order.DiscountAmount,
			Total:          order.Subtotal + order.Tax + order.DeliveryFee - order.DiscountAmount,
			Status:         order.Status,
			CreatedAt:      order.CreatedAt,
			Items:          order.Items,
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

	err = s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus)
	if err != nil {
		return err
	}

	// Insert tracking log
	desc := fmt.Sprintf("Order status updated to %s", newStatus)
	_ = s.trackingRepo.InsertTracking(ctx, nil, orderID, newStatus, &desc)

	// Earn points if the order becomes paid/completed/delivered
	// Just an example: 1 point for every 1000 IDR subtotal
	if newStatus == "paid" || newStatus == "completed" || newStatus == "delivered" {
		if order.Status != "paid" && order.Status != "completed" && order.Status != "delivered" {
			earnedPoints := int(order.Subtotal / 1000)
			if earnedPoints > 0 {
				_ = s.pointService.AddPoints(ctx, nil, customerID, &orderID, earnedPoints, "Earned points from order completion")
			}
		}
	}

	return nil
}

func (s *OrderService) GetOrderTracking(ctx context.Context, orderID int, customerID int) ([]models.OrderTracking, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.CustomerID != customerID {
		return nil, fmt.Errorf("unauthorized access to order")
	}

	return s.trackingRepo.GetTrackingByOrderID(ctx, orderID)
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

func (s *OrderService) GetOrderStats(ctx context.Context) (map[string]int, error) {
	stats, err := s.orderRepo.GetOrderStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order stats: %w", err)
	}
	return stats, nil
}

func (s *OrderService) GenerateSnapToken(orderID int, grossAmount float64) (string, string, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = "test-server-key"
	}

	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"

	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     fmt.Sprintf("%d", orderID),
			"gross_amount": grossAmount,
		},
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(serverKey+":"))
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token         string   `json:"token"`
		RedirectURL   string   `json:"redirect_url"`
		ErrorMessages []string `json:"error_messages"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if len(result.ErrorMessages) > 0 {
		return "", "", fmt.Errorf("midtrans error: %v", result.ErrorMessages)
	}

	return result.Token, result.RedirectURL, nil
}