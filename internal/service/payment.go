package service

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"os"
	"strconv"
	"strings"

	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type PaymentService struct {
	paymentRepo  *repository.PaymentRepository
	orderRepo    *repository.OrderRepository
	pointService PointService
	trackingRepo repository.TrackingRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository, pointService PointService, trackingRepo repository.TrackingRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo, pointService: pointService, trackingRepo: trackingRepo}
}

func (s *PaymentService) Create(ctx context.Context, payment *models.Payment) error {
	if payment.OrderID <= 0 {
		return errors.NewValidationError("order_id", "must be positive")
	}
	if payment.Amount <= 0 {
		return errors.NewValidationError("amount", "must be positive")
	}

	validMethods := map[string]bool{"credit_card": true, "debit_card": true, "bank_transfer": true, "e_wallet": true}
	if !validMethods[payment.Method] {
		return errors.NewValidationError("method", "invalid payment method")
	}

	// Verify order exists
	_, err := s.orderRepo.GetOrderByID(ctx, payment.OrderID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return errors.NewValidationError("order_id", "order not found")
		}
		return errors.NewServiceError("verify order", err)
	}

	return s.paymentRepo.Create(ctx, payment)
}

func (s *PaymentService) GetByID(ctx context.Context, id int) (*models.Payment, error) {
	if id <= 0 {
		return nil, errors.NewValidationError("id", "must be positive")
	}

	payment, err := s.paymentRepo.GetByOrderID(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, err
		}
		return nil, errors.NewServiceError("get payment", err)
	}
	return payment, nil
}

func (s *PaymentService) GetByOrderID(ctx context.Context, orderID int) (*models.Payment, error) {
	if orderID <= 0 {
		return nil, errors.NewValidationError("order_id", "must be positive")
	}

	payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, err
		}
		return nil, errors.NewServiceError("get payment", err)
	}
	return payment, nil
}

func (s *PaymentService) UpdateStatus(ctx context.Context, id int, status string) error {
	if id <= 0 {
		return errors.NewValidationError("id", "must be positive")
	}

	validStatuses := map[string]bool{"pending": true, "completed": true, "failed": true, "cancelled": true}
	if !validStatuses[status] {
		return errors.NewValidationError("status", "invalid payment status")
	}

	err := s.paymentRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return err
		}
		return errors.NewServiceError("update payment status", err)
	}
	return nil
}

func (s *PaymentService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.NewValidationError("id", "must be positive")
	}

	err := s.paymentRepo.Delete(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return err
		}
		return errors.NewServiceError("delete payment", err)
	}
	return nil
}

func (s *PaymentService) VerifyMidtransSignature(orderID, statusCode, grossAmount, signatureKey, serverKey string) bool {
	payload := orderID + statusCode + grossAmount + serverKey
	hasher := sha512.New()
	hasher.Write([]byte(payload))
	computedHash := hex.EncodeToString(hasher.Sum(nil))
	return strings.ToLower(computedHash) == strings.ToLower(signatureKey)
}

func (s *PaymentService) ProcessMidtransCallback(ctx context.Context, req models.MidtransCallbackRequest) error {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		serverKey = "test-server-key"
	}

	if !s.VerifyMidtransSignature(req.OrderID, req.StatusCode, req.GrossAmount, req.SignatureKey, serverKey) {
		return errors.NewValidationError("signature_key", "invalid signature key")
	}

	orderID, err := strconv.Atoi(req.OrderID)
	if err != nil {
		return errors.NewValidationError("order_id", "invalid order_id format")
	}

	payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return errors.NewValidationError("order_id", "payment not found for order")
		}
		return errors.NewServiceError("get payment", err)
	}

	if payment.Status == "completed" || payment.Status == "failed" || payment.Status == "cancelled" {
		return nil // Idempotent
	}

	grossAmount, _ := strconv.ParseFloat(req.GrossAmount, 64)

	var newPaymentStatus string
	var newOrderStatus string
	var trackingMsg string

	if req.TransactionStatus == "settlement" || req.TransactionStatus == "capture" {
		if req.FraudStatus == "challenge" {
			newPaymentStatus = "pending"
			newOrderStatus = "pending"
			trackingMsg = "Payment challenged by Midtrans fraud detection"
		} else {
			newPaymentStatus = "completed"
			newOrderStatus = "processing"
			trackingMsg = "Payment settled successfully via Midtrans"

			pointsToReward := int(grossAmount / 10000)
			
			order, _ := s.orderRepo.GetOrderByID(ctx, orderID)
			if order != nil && pointsToReward > 0 {
				_ = s.pointService.AddPoints(ctx, nil, order.CustomerID, &orderID, pointsToReward, "Reward from order payment")
			}
		}
	} else if req.TransactionStatus == "cancel" || req.TransactionStatus == "expire" || req.TransactionStatus == "deny" {
		newPaymentStatus = "failed"
		newOrderStatus = "cancelled"
		trackingMsg = "Payment failed or expired"
	} else {
		return nil
	}

	if newPaymentStatus != "" && newPaymentStatus != payment.Status {
		err = s.paymentRepo.UpdateStatus(ctx, payment.ID, newPaymentStatus)
		if err != nil {
			return errors.NewServiceError("update payment status", err)
		}
	}

	if newOrderStatus != "" {
		err = s.orderRepo.UpdateOrderStatus(ctx, orderID, newOrderStatus)
		if err != nil {
			return errors.NewServiceError("update order status", err)
		}
	}

	if trackingMsg != "" {
		_ = s.trackingRepo.InsertTracking(ctx, nil, orderID, newOrderStatus, &trackingMsg)
	}

	return nil
}
