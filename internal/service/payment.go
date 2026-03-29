package service

import (
	"context"
	"koda-b6-backend/internal/errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo}
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

	payment, err := s.paymentRepo.GetByID(ctx, id)
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
