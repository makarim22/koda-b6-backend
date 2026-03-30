package service

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"
	"koda-b6-backend/internal/repository"
)

type ProductDiscountService struct {
	repo *repository.ProductDiscountRepository
}

func NewProductDiscountService(repo *repository.ProductDiscountRepository) *ProductDiscountService {
	return &ProductDiscountService{
		repo: repo,
	}
}

func (s *ProductDiscountService) GetDiscountsByProductID(ctx context.Context, productID int) ([]models.ProductDiscount, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("invalid product ID: %d", productID)
	}
	discounts, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch variants for product: %w", err)
	}
	return discounts, nil
}
