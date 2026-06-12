package service

import (
	"context"
	"errors"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
	"time"
)

type VoucherService struct {
	voucherRepo *repository.VoucherRepository
}

func NewVoucherService(voucherRepo *repository.VoucherRepository) *VoucherService {
	return &VoucherService{voucherRepo: voucherRepo}
}

// CalculateDiscount validates the voucher and computes the discount
func (s *VoucherService) CalculateDiscount(ctx context.Context, code string, subtotal float64) (*models.Voucher, float64, error) {
	voucher, err := s.voucherRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, 0, errors.New("invalid voucher code")
	}

	// 1. Check validity dates
	now := time.Now()
	if now.Before(voucher.ValidFrom) {
		return nil, 0, errors.New("voucher is not yet active")
	}
	if now.After(voucher.ValidUntil) {
		return nil, 0, errors.New("voucher has expired")
	}

	// 2. Check usage limit
	if voucher.UsageLimit > 0 && voucher.UsedCount >= voucher.UsageLimit {
		return nil, 0, errors.New("voucher usage limit reached")
	}

	// 3. Check minimum purchase
	if subtotal < voucher.MinPurchase {
		return nil, 0, errors.New("minimum purchase requirement not met")
	}

	// 4. Calculate discount amount
	var discountAmount float64
	if voucher.DiscountType == "PERCENTAGE" {
		discountAmount = subtotal * (voucher.DiscountValue / 100)
		if voucher.MaxDiscount != nil && discountAmount > *voucher.MaxDiscount {
			discountAmount = *voucher.MaxDiscount
		}
	} else if voucher.DiscountType == "FIXED" {
		discountAmount = voucher.DiscountValue
	}

	// Prevent discount from exceeding subtotal
	if discountAmount > subtotal {
		discountAmount = subtotal
	}

	return voucher, discountAmount, nil
}

func (s *VoucherService) GetAll(ctx context.Context) ([]models.Voucher, error) {
	return s.voucherRepo.GetAll(ctx)
}

func (s *VoucherService) Create(ctx context.Context, voucher *models.Voucher) error {
	return s.voucherRepo.Create(ctx, voucher)
}

func (s *VoucherService) Update(ctx context.Context, id int, voucher *models.Voucher) error {
	return s.voucherRepo.Update(ctx, id, voucher)
}

func (s *VoucherService) Delete(ctx context.Context, id int) error {
	return s.voucherRepo.Delete(ctx, id)
}
